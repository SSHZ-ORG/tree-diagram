package models

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type EventSnapshot struct {
	EventID   string    `datastore:",noindex" json:"-"` // Use Key.Parent
	Timestamp time.Time `json:"timestamp"`

	// Recorded even if not changed.
	NoteCount int `datastore:",noindex" json:"noteCount"` // Don't query Snapshots directly.

	// Recorded only if changed.
	Actors []*datastore.Key `datastore:",noindex" json:"-"`
}

const eventSnapshotKind = "EventSnapshot"

func maybeCreateSnapshot(ctx context.Context, ek *datastore.Key, oe, ne *Event, ts time.Time) (*datastore.Key, *EventSnapshot) {
	s := &EventSnapshot{
		EventID:   ne.ID,
		Timestamp: ts,
		NoteCount: ne.LastNoteCount,
	}

	if oe == nil {
		oe = &Event{}
	}

	if !areKeysSetsEqual(oe.Actors, ne.Actors) {
		s.Actors = ne.Actors
	}

	log.Debugf(ctx, "Taking snapshot for event %s. (%d -> %d)", ne.debugName(), oe.LastNoteCount, ne.LastNoteCount)
	return datastore.NewIncompleteKey(ctx, eventSnapshotKind, ek), s
}

// Errors wrapped.
// This bypasses memcache.
func getSnapshotsForEvent(ctx context.Context, eventKey *datastore.Key) ([]*EventSnapshot, error) {
	var es []*EventSnapshot
	_, err := datastore.NewQuery(eventSnapshotKind).Ancestor(eventKey).Order("Timestamp").GetAll(ctx, &es)
	if err != nil {
		return nil, errors.Wrap(err, "datastore query failed")
	}

	return es, nil
}
