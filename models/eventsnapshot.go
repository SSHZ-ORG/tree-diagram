package models

import (
	"time"

	"github.com/pkg/errors"
	"github.com/qedus/nds"
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
func getSnapshotsForEvent(ctx context.Context, eventKey *datastore.Key) ([]*EventSnapshot, error) {
	keys, err := datastore.NewQuery(eventSnapshotKind).Ancestor(eventKey).Order("Timestamp").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "datastore query failed")
	}

	es := make([]*EventSnapshot, len(keys))
	err = nds.GetMulti(ctx, keys, es)
	if err != nil {
		return nil, errors.Wrap(err, "nds.GetMulti failed")
	}

	return es, nil
}
