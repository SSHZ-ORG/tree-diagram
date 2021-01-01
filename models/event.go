package models

import (
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/models/cache"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/scylladb/go-set/strset"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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

type FrontendEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Date     string `json:"date"`
	Finished bool   `json:"finished"`

	LastNoteCount int `json:"lastNoteCount"`
}

func (e Event) debugName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Name)
}

func (e Event) ToFrontendEvent() FrontendEvent {
	return FrontendEvent{
		ID:   e.ID,
		Name: e.Name,

		Date:     civil.DateOf(e.Date).String(),
		Finished: e.Finished,

		LastNoteCount: e.LastNoteCount,
	}
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
func QueryEvents(ctx context.Context, placeID string, actorIDs []string, limit, offset int) ([]*Event, error) {
	query := datastore.NewQuery(eventKind).KeysOnly().Limit(limit).Offset(offset).Order("-LastNoteCount")

	if placeID != "" {
		query = query.Filter("Place =", getPlaceKey(ctx, placeID))
	}

	for _, actorID := range actorIDs {
		query = query.Filter("Actors =", getActorKey(ctx, actorID))
	}

	keys, err := query.GetAll(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "query.GetAll failed")
	}

	es := make([]*Event, len(keys))
	err = nds.GetMulti(ctx, keys, es)

	return es, errors.Wrap(err, "nds.GetMulti failed")
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
	response.Date = proto.String(civil.DateOf(e.Date).String())

	var errTotal, errFinished error
	wg := sync.WaitGroup{}
	wg.Add(2)
	baseQuery := datastore.NewQuery(eventKind).Filter("Place =", e.Place)
	finishedQuery := baseQuery.Filter("Finished =", true)

	go func() {
		defer wg.Done()
		if t, err := baseQuery.Count(ctx); err == nil {
			response.PlaceStatsTotal.Total = proto.Int32(int32(t))
		} else {
			errTotal = errors.Wrap(err, "Count total total failed")
			return
		}
		if r, err := baseQuery.Filter("LastNoteCount >", e.LastNoteCount).Count(ctx); err == nil {
			response.PlaceStatsTotal.Rank = proto.Int32(int32(r))
		} else {
			errTotal = errors.Wrap(err, "Count total rank failed")
			return
		}
	}()
	go func() {
		defer wg.Done()
		if t, err := finishedQuery.Count(ctx); err == nil {
			response.PlaceStatsFinished.Total = proto.Int32(int32(t))
		} else {
			errFinished = errors.Wrap(err, "Count finished total failed")
			return
		}
		if r, err := finishedQuery.Filter("LastNoteCount >", e.LastNoteCount).Count(ctx); err == nil {
			response.PlaceStatsFinished.Rank = proto.Int32(int32(r))
		} else {
			errFinished = errors.Wrap(err, "Count finished rank failed")
			return
		}
	}()

	snapshots, err := getSnapshotsForEvent(ctx, key)
	if err != nil {
		return nil, err
	}

	akSet := strset.New()
	for _, s := range snapshots {
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

	actorNames := make(map[string]string)
	actors, err := getActors(ctx, aks)
	if err != nil {
		return response, err
	}
	for i, a := range actors {
		actorNames[aks[i].Encode()] = a.Name
	}

	lastActors := strset.New()
	for _, s := range snapshots {
		item := &pb.RenderEventResponse_Snapshot{
			Timestamp: timestamppb.New(s.Timestamp),
			NoteCount: proto.Int32(int32(s.NoteCount)),
		}

		if len(s.Actors) > 0 {
			newActors := strset.New()
			for _, ak := range s.Actors {
				newActors.Add(ak.Encode())
			}

			strset.Difference(newActors, lastActors).Each(func(addedKey string) bool {
				item.AddedActors = append(item.AddedActors, actorNames[addedKey])
				return true
			})
			strset.Difference(lastActors, newActors).Each(func(removedKey string) bool {
				item.RemovedActors = append(item.RemovedActors, actorNames[removedKey])
				return true
			})
			lastActors = newActors
		}

		response.Snapshots = append(response.Snapshots, item)
	}

	wg.Wait()
	if errTotal != nil {
		return nil, errTotal
	}
	if errFinished != nil {
		return nil, errFinished
	}

	return response, nil
}
