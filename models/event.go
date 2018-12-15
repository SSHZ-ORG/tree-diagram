package models

import (
	"time"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Event struct {
	ID   string
	Name string

	// Used as civil date. Always at UTC midnight.
	Date time.Time

	Place  *datastore.Key
	Actors []*datastore.Key

	// These 3 are stored as local time (JST).
	OpenTime  time.Time
	StartTime time.Time
	EndTime   time.Time

	LastNoteCount int

	LastUpdateTime time.Time
}

const eventKind = "Event"

// Insert an Event. This automatically takes an snapshot.
func InsertOrUpdateEvent(ctx context.Context, e *Event) error {
	key := datastore.NewKey(ctx, eventKind, e.ID, 0, nil)

	return nds.RunInTransaction(ctx, func(tc context.Context) error {
		_, err := nds.Put(ctx, key, e)
		if err != nil {
			return err
		}

		return TakeSnapshot(ctx, key, e)
	}, nil)
}
