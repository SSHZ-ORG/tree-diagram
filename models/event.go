package models

import (
	"fmt"
	"time"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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

func (e Event) debugName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Name)
}

const eventKind = "Event"

// Insert an Event. This automatically takes an snapshot.
func InsertOrUpdateEvent(ctx context.Context, e *Event) error {
	key := datastore.NewKey(ctx, eventKind, e.ID, 0, nil)

	oe := &Event{}
	err := nds.Get(ctx, key, oe)
	if err != nil && err != datastore.ErrNoSuchEntity {
		// Something else happens. Rethrow it.
		return err
	}

	// We always put the event even if nothing changed, to update the timestamp.
	// But a snapshot is taken only if needSnapshot says so.
	return nds.RunInTransaction(ctx, func(tc context.Context) error {
		_, err := nds.Put(ctx, key, e)
		if err != nil {
			return err
		}

		if needSnapshot(oe, e) {
			log.Debugf(ctx, "Taking snapshot for event %s.", e.debugName())
			return takeSnapshot(ctx, key, e)
		} else {
			log.Debugf(ctx, "Snapshot skipped for event %s.", e.debugName())
			return nil
		}
	}, nil)
}
