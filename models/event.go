package models

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/qedus/nds"
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

	LastUpdateTime time.Time
}

type FrontendEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Date string `json:"date"`

	LastNoteCount int `json:"last_note_count"`
}

func (e Event) debugName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Name)
}

func (e Event) ToFrontendEvent() FrontendEvent {
	return FrontendEvent{
		ID:   e.ID,
		Name: e.Name,

		Date: civil.DateOf(e.Date).String(),

		LastNoteCount: e.LastNoteCount,
	}
}

const (
	eventKind     = "Event"
	queryPageSize = 30
)

func getEventKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, eventKind, id, 0, nil)
}

// Insert or update events. This automatically takes snapshots if needed.
func InsertOrUpdateEvents(ctx context.Context, events []*Event) error {
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
					return err
				}
			}
		} else {
			return err
		}
	}

	var snapshotKeys []*datastore.Key
	var snapshots []*EventSnapshot
	for i, e := range events {
		sk, s := maybeCreateSnapshot(ctx, keys[i], oes[i], e)
		if sk != nil {
			snapshotKeys = append(snapshotKeys, sk)
			snapshots = append(snapshots, s)
		}
	}

	// We always update events even if no change, to update the timestamp.
	if _, err := nds.PutMulti(ctx, keys, events); err != nil {
		return err
	}
	if _, err := nds.PutMulti(ctx, snapshotKeys, snapshots); err != nil {
		return err
	}

	return nil
}

func QueryEvents(ctx context.Context, placeID string, actorIDs []string, page int) ([]*Event, error) {
	query := datastore.NewQuery(eventKind).KeysOnly().Limit(queryPageSize).Order("-LastNoteCount")

	if placeID != "" {
		query = query.Filter("Place=", getPlaceKey(ctx, placeID))
	}

	for _, actorID := range actorIDs {
		query = query.Filter("Actors=", getActorKey(ctx, actorID))
	}

	if page > 1 {
		query = query.Offset(queryPageSize * (page - 1))
	}

	keys, err := query.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	es := make([]*Event, len(keys))
	err = nds.GetMulti(ctx, keys, es)

	return es, err
}

type RenderEventResponse struct {
	Date      string           `json:"date"`
	Snapshots []*EventSnapshot `json:"snapshots"`
}

func PrepareRenderEventResponse(ctx context.Context, eventID string) (*RenderEventResponse, error) {
	key := getEventKey(ctx, eventID)

	response := &RenderEventResponse{
		Snapshots: make([]*EventSnapshot, 0), // So json does not make it null.
	}

	e := &Event{}
	if err := nds.Get(ctx, key, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return response, nil // Don't care if we don't know about the event yet.
		}
		return nil, err
	}
	response.Date = civil.DateOf(e.Date).String()

	snapshots, err := getSnapshotsForEvent(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(snapshots) > 0 {
		response.Snapshots = snapshots
	}

	return response, nil
}
