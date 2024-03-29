package models

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/models/cache"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/scylladb/go-set/strset"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/datastore"
	"google.golang.org/appengine/v2/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Event struct {
	ID   string `datastore:",noindex"` // Use Key
	Name string

	Date     time.Time // Used as civil date. Always at UTC midnight.
	Finished bool      // Use this as Date < Now().

	Place  *datastore.Key
	Actors []*datastore.Key

	// These 3 are stored as local time (JST).
	OpenTime  time.Time
	StartTime time.Time
	EndTime   time.Time

	LastNoteCount int

	// The last time we see something changed in the event. Does not affect Equals().
	LastUpdateTime time.Time
}

func (e *Event) debugName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Name)
}

const (
	eventKind = "Event"
)

func getEventKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, eventKind, id, 0, nil)
}

// Insert or update events. This automatically takes snapshots if needed.
// Returns actors whose related events have changed (and thus apicache should be cleared.)
// Errors wrapped.
func InsertOrUpdateEvents(ctx context.Context, events []*Event) (*strset.Set, error) {
	if len(events) == 0 {
		return strset.New(), nil
	}

	batchSize := len(events)
	if batchSize >= minEventsToParallelize {
		batchSize = (len(events) + 1) / 2
	}
	if batchSize > maxEntitiesPerXGTransaction {
		batchSize = maxEntitiesPerXGTransaction
	}

	var batches [][]*Event
	for batchSize < len(events) {
		events, batches = events[batchSize:], append(batches, events[0:batchSize:batchSize])
	}
	batches = append(batches, events)

	wg := sync.WaitGroup{}
	wg.Add(len(batches))
	errs := make([]error, len(batches))
	diffActorIDs := make([]*strset.Set, len(batches))

	for i := range batches {
		go func(i int) {
			defer wg.Done()
			diffActorIDs[i], errs[i] = internalInsertOrUpdateEvents(ctx, batches[i])
		}(i)
	}

	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return strset.Union(diffActorIDs...), nil
}

// Insert or update at most 25 events because of Transaction XG restriction.
// We can remove the restriction after we are migrated to Firestore.
// Returns actors whose related events have changed (and thus apicache should be cleared.)
// Errors wrapped.
func internalInsertOrUpdateEvents(ctx context.Context, events []*Event) (*strset.Set, error) {
	if len(events) == 0 {
		return nil, nil
	}

	var eventKeys []*datastore.Key
	for _, e := range events {
		eventKeys = append(eventKeys, getEventKey(ctx, e.ID))
	}

	oes := make([]*Event, len(events))
	err := nds.GetMulti(ctx, eventKeys, oes)
	if err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for _, e := range me {
				if e != nil && e != datastore.ErrNoSuchEntity {
					// Something else happened. Rethrow it.
					return nil, errors.Wrap(err, "nds.GetMulti returned error other than NoSuchEntity")
				}
			}
		} else {
			return nil, errors.Wrap(err, "nds.GetMulti returned error that is not a MultiError")
		}
	}

	var insertedCESKeys []*datastore.Key
	diffActorIDs := strset.New()

	err = nds.RunInTransaction(ctx, func(ctx context.Context) error {
		var keysToInsert []*datastore.Key
		var eventsToInsert []*Event

		var cesKeysToInsert []*datastore.Key
		var cessToInsert []*compressedEventSnapshot

		var cesKeysToAppend []*datastore.Key
		cesToAppendIndexes := make(map[int]int) // Map from index in events -> index in cesKeysToAppend

		for i, e := range events {
			if !e.Equal(oes[i]) {
				// Need to update the event itself.
				keysToInsert = append(keysToInsert, eventKeys[i])
				eventsToInsert = append(eventsToInsert, e)
				diffActorIDs = strset.Union(diffActorIDs, e.DiffActorIDs(oes[i]))
			}

			cesKey, err := maybeGetCESKeyToAppend(ctx, oes[i], e, eventKeys[i])
			if err != nil {
				return err
			}

			if cesKey != nil {
				// We should append to an existing CES. Prepare to fetch it.
				cesKeysToAppend = append(cesKeysToAppend, cesKey)
				cesToAppendIndexes[i] = len(cesKeysToAppend) - 1
			}
		}

		cessToAppend := make([]*compressedEventSnapshot, len(cesKeysToAppend))
		if err := nds.GetMulti(ctx, cesKeysToAppend, cessToAppend); err != nil {
			return errors.Wrap(err, "nds.GetMulti failed")
		}

		for i, e := range events {
			var cesKey *datastore.Key
			var ces *compressedEventSnapshot

			if idx, ok := cesToAppendIndexes[i]; ok {
				// We should append to an existing CES.
				cesKey, ces = cesKeysToAppend[idx], cessToAppend[idx]

				if ces.isConsistent(e) {
					// Good. Let's append.
					ces.appendTime(e.LastUpdateTime)
					log.Debugf(ctx, "Appending to CES %+v for event %s (%d -> %d)", cesKey, e.debugName(), oes[i].LastNoteCount, e.LastNoteCount)
				} else {
					// Something is wrong. Shout out and create new CES.
					cesKey, ces = nil, nil
					log.Criticalf(ctx, "Inconsistent CES %+v for event %s detected!", cesKey, e.debugName())
				}
			}

			if cesKey == nil {
				// Either maybeGetCESKeyToAppend thinks we should create new one, or the last CES is inconsistent.
				cesKey, ces = newCESFromEvent(ctx, oes[i], e, eventKeys[i])

				lastNoteCount := 0
				if oes[i] != nil {
					lastNoteCount = oes[i].LastNoteCount
				}
				log.Debugf(ctx, "Creating new CES for event %s (%d -> %d)", e.debugName(), lastNoteCount, e.LastNoteCount)
			}

			cesKeysToInsert = append(cesKeysToInsert, cesKey)
			cessToInsert = append(cessToInsert, ces)
		}

		if err := cache.ClearLastCESKeys(ctx, eventKeys); err != nil {
			return err
		}

		if _, err := nds.PutMulti(ctx, keysToInsert, eventsToInsert); err != nil {
			return errors.Wrap(err, "nds.PutMulti failed")
		}

		insertedCESKeys, err = nds.PutMulti(ctx, cesKeysToInsert, cessToInsert)
		if err != nil {
			return errors.Wrap(err, "nds.PutMulti failed")
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		return nil, err
	}

	cache.SetLastCESKeys(ctx, insertedCESKeys)
	return diffActorIDs, nil
}

// This ignores LastUpdateTime
func (e *Event) Equal(o *Event) bool {
	if e != nil && o != nil {
		if e.ID != o.ID || e.Name != o.Name || e.Date != o.Date || e.Finished != o.Finished {
			return false
		}
		if !e.Place.Equal(o.Place) {
			return false
		}
		if !areKeysSetsEqual(e.Actors, o.Actors) {
			return false
		}
		if !e.OpenTime.Equal(o.OpenTime) || !e.StartTime.Equal(o.StartTime) || !e.EndTime.Equal(o.EndTime) {
			return false
		}
		if e.LastNoteCount != o.LastNoteCount {
			return false
		}
		return true
	}
	return e == o
}

func (e *Event) DiffActorIDs(o *Event) *strset.Set {
	as := strset.New()
	bs := strset.New()

	if e != nil {
		for _, k := range e.Actors {
			as.Add(k.StringID())
		}
	}
	if o != nil {
		for _, k := range o.Actors {
			bs.Add(k.StringID())
		}
	}

	return strset.SymmetricDifference(as, bs)
}

// Set Finished to False for everything that happened before today.
// For events on yesterday, the daily cron is not executed yet when this cron happens.
// Not using Transaction as this should never overlap with daily cron.
func CleanupFinishedEvents(ctx context.Context, today civil.Date) error {
	query := datastore.NewQuery(eventKind).KeysOnly().Filter("Date <", today.In(time.UTC)).Filter("Finished =", false)

	keys, err := query.GetAll(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "query.GetAll failed")
	}

	es := make([]*Event, len(keys))
	err = nds.GetMulti(ctx, keys, es)
	if err != nil {
		return errors.Wrap(err, "nds.GetMulti failed")
	}

	for _, e := range es {
		e.Finished = true
	}

	_, err = nds.PutMulti(ctx, keys, es)
	return errors.Wrap(err, "nds.PutMulti failed")
}

// Errors wrapped.
func QueryEvents(ctx context.Context, filter *pb.QueryEventsRequest_EventFilter, limit, offset int) ([]*Event, bool, error) {
	query := datastore.NewQuery(eventKind).KeysOnly().Limit(limit + 1).Offset(offset).Order("-LastNoteCount")

	if filter.GetPlaceId() != "" {
		query = query.Filter("Place =", getPlaceKey(ctx, filter.GetPlaceId()))
	}

	for _, actorID := range filter.GetActorIds() {
		query = query.Filter("Actors =", getActorKey(ctx, actorID))
	}

	keys, err := query.GetAll(ctx, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "query.GetAll failed")
	}

	hasNext := false
	if len(keys) > limit {
		hasNext = true
		keys = keys[:limit]
	}

	es := make([]*Event, len(keys))
	err = nds.GetMulti(ctx, keys, es)

	return es, hasNext, errors.Wrap(err, "nds.GetMulti failed")
}

// Errors wrapped.
func PrepareRenderEventResponse(ctx context.Context, eventID string) (*pb.RenderEventResponse, error) {
	key := getEventKey(ctx, eventID)

	response := &pb.RenderEventResponse{
		PlaceStatsTotal:    &pb.RenderEventResponse_PlaceNoteCountStats{},
		PlaceStatsFinished: &pb.RenderEventResponse_PlaceNoteCountStats{},
	}

	e := &Event{}
	if err := nds.Get(ctx, key, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return response, nil // Don't care if we don't know about the event yet.
		}
		return nil, errors.Wrap(err, "nds.Get failed")
	}
	response.Date = utils.ToProtoDate(civil.DateOf(e.Date))

	g, gctx := errgroup.WithContext(ctx)
	baseQuery := datastore.NewQuery(eventKind).Filter("Place =", e.Place)
	finishedQuery := baseQuery.Filter("Finished =", true)
	g.Go(func() error {
		t, err := baseQuery.Count(gctx)
		response.PlaceStatsTotal.Total = proto.Int32(int32(t))
		return errors.Wrap(err, "Count total total failed")
	})
	g.Go(func() error {
		r, err := baseQuery.Filter("LastNoteCount >", e.LastNoteCount).Count(gctx)
		response.PlaceStatsTotal.Rank = proto.Int32(int32(r))
		return errors.Wrap(err, "Count total rank failed")
	})
	g.Go(func() error {
		t, err := finishedQuery.Count(gctx)
		response.PlaceStatsFinished.Total = proto.Int32(int32(t))
		return errors.Wrap(err, "Count finished total failed")
	})
	g.Go(func() error {
		r, err := finishedQuery.Filter("LastNoteCount >", e.LastNoteCount).Count(gctx)
		response.PlaceStatsFinished.Rank = proto.Int32(int32(r))
		return errors.Wrap(err, "Count finished rank failed")
	})

	compressedSnapshots, err := getCompressedSnapshots(ctx, key)
	if err != nil {
		return nil, err
	}

	akSet := strset.New()
	for _, s := range compressedSnapshots {
		for _, ak := range s.Actors {
			akSet.Add(ak.Encode())
		}
	}
	var aks []*datastore.Key
	akSet.Each(func(i string) bool {
		dk, _ := datastore.DecodeKey(i)
		aks = append(aks, dk)
		return true
	})

	actorInfos := make(map[string]*pb.RenderEventResponse_ActorInfo)
	actors, err := getActors(ctx, aks)
	if err != nil {
		return response, err
	}
	for i, a := range actors {
		actorInfos[aks[i].Encode()] = &pb.RenderEventResponse_ActorInfo{
			Id:   &a.ID,
			Name: &a.Name,
		}
	}

	var lastActors []string
	las := strset.New()
	for _, s := range compressedSnapshots {
		item := &pb.RenderEventResponse_CompressedSnapshot{
			NoteCount: proto.Int32(int32(s.NoteCount)),
		}
		for _, ts := range s.Timestamps {
			item.Timestamps = append(item.Timestamps, timestamppb.New(ts))
		}

		if len(s.Actors) > 0 {
			var newActors []string
			for _, ak := range s.Actors {
				newActors = append(newActors, ak.Encode())
			}

			for _, a := range newActors {
				if !las.Has(a) {
					item.AddedActors = append(item.AddedActors, actorInfos[a])
				}
			}

			nas := strset.New(newActors...)
			for _, a := range lastActors {
				if !nas.Has(a) {
					item.RemovedActors = append(item.RemovedActors, actorInfos[a])
				}
			}

			lastActors = newActors
			las = nas
		}

		response.CompressedSnapshots = append(response.CompressedSnapshots, item)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return response, nil
}
