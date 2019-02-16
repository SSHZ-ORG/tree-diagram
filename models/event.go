package models

import (
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/civil"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/scylladb/go-set/strset"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
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
// Errors wrapped.
func InsertOrUpdateEvents(ctx context.Context, events []*Event, ts time.Time) error {
	if len(events) == 0 {
		return nil
	}

	var keys []*datastore.Key
	for _, e := range events {
		keys = append(keys, getEventKey(ctx, e.ID))
	}

	oes := make([]*Event, len(events))
	err := nds.GetMulti(ctx, keys, oes)
	if err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for _, e := range me {
				if e != nil && e != datastore.ErrNoSuchEntity {
					// Something else happened. Rethrow it.
					return errors.Wrap(err, "nds.GetMulti returned error other than NoSuchEntity")
				}
			}
		} else {
			return errors.Wrap(err, "nds.GetMulti returned error that is not a MultiError")
		}
	}

	var snapshotKeys []*datastore.Key
	var snapshots []*EventSnapshot
	var keysToInsert []*datastore.Key
	var eventsToInsert []*Event
	for i, e := range events {
		sk, s := maybeCreateSnapshot(ctx, keys[i], oes[i], e, ts)
		if sk != nil {
			snapshotKeys = append(snapshotKeys, sk)
			snapshots = append(snapshots, s)
		}

		if !e.Equal(oes[i]) {
			e.LastUpdateTime = ts
			keysToInsert = append(keysToInsert, keys[i])
			eventsToInsert = append(eventsToInsert, e)
		}
	}

	if _, err := nds.PutMulti(ctx, keysToInsert, eventsToInsert); err != nil {
		return errors.Wrap(err, "nds.PutMulti failed")
	}
	if _, err := nds.PutMulti(ctx, snapshotKeys, snapshots); err != nil {
		return errors.Wrap(err, "nds.PutMulti failed")
	}

	return nil
}

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

// If an event was picked up by us but disappear later (deleted / de-duped / NoteCount fell below threshold),
// its Finished won't be updated. Clean them up manually.
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

type PlaceNoteCountStats struct {
	Total int `json:"total"`
	Rank  int `json:"rank"`
}

type FrontendEventSnapshot struct {
	Timestamp time.Time `json:"timestamp"`
	NoteCount int       `json:"noteCount"`

	AddedActors   []string `json:"addedActors"`
	RemovedActors []string `json:"removedActors"`
}

type RenderEventResponse struct {
	Date      string                   `json:"date"`
	Snapshots []*FrontendEventSnapshot `json:"snapshots"`

	PlaceStatsTotal    PlaceNoteCountStats `json:"placeStatsTotal"`
	PlaceStatsFinished PlaceNoteCountStats `json:"placeStatsFinished"`
}

// Errors wrapped.
func PrepareRenderEventResponse(ctx context.Context, eventID string) (*RenderEventResponse, error) {
	key := getEventKey(ctx, eventID)

	response := &RenderEventResponse{
		Snapshots: make([]*FrontendEventSnapshot, 0), // So json does not make it null.
	}

	e := &Event{}
	if err := nds.Get(ctx, key, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return response, nil // Don't care if we don't know about the event yet.
		}
		return nil, errors.Wrap(err, "nds.Get failed")
	}
	response.Date = civil.DateOf(e.Date).String()

	var errTotal, errFinished error
	wg := sync.WaitGroup{}
	wg.Add(2)
	baseQuery := datastore.NewQuery(eventKind).Filter("Place =", e.Place)
	finishedQuery := baseQuery.Filter("Finished =", true)

	go func() {
		defer wg.Done()
		var err error
		if response.PlaceStatsTotal.Total, err = baseQuery.Count(ctx); err != nil {
			errTotal = errors.Wrap(err, "Count total total failed")
			return
		}
		if response.PlaceStatsTotal.Rank, err = baseQuery.Filter("LastNoteCount >", e.LastNoteCount).Count(ctx); err != nil {
			errTotal = errors.Wrap(err, "Count total rank failed")
			return
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		if response.PlaceStatsFinished.Total, err = finishedQuery.Count(ctx); err != nil {
			errFinished = errors.Wrap(err, "Count finished total failed")
			return
		}
		if response.PlaceStatsFinished.Rank, err = finishedQuery.Filter("LastNoteCount >", e.LastNoteCount).Count(ctx); err != nil {
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
	actors, err := GetActors(ctx, aks)
	if err != nil {
		return response, err
	}
	for i, a := range actors {
		actorNames[aks[i].Encode()] = a.Name
	}

	lastActors := strset.New()
	for _, s := range snapshots {
		item := &FrontendEventSnapshot{
			Timestamp:     s.Timestamp,
			NoteCount:     s.NoteCount,
			AddedActors:   make([]string, 0),
			RemovedActors: make([]string, 0),
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
