package models

import (
	"time"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type EventSnapshot struct {
	EventID   string    `datastore:",noindex" json:"-"` // Use Key.Parent
	Timestamp time.Time `json:"timestamp"`

	// Recorded even if not changed.
	NoteCount int `datastore:",noindex" json:"note_count"` // Don't query Snapshots directly.

	// Recorded only if changed.
	Actors []*datastore.Key `datastore:",noindex" json:"-"`
}

const eventSnapshotKind = "EventSnapshot"

func maybeCreateSnapshot(ctx context.Context, ek *datastore.Key, oe, ne *Event) (*datastore.Key, *EventSnapshot) {
	s := &EventSnapshot{
		EventID:   ne.ID,
		Timestamp: ne.LastUpdateTime,
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

func GetSnapshotsForEvent(ctx context.Context, eventID string) ([]*EventSnapshot, error) {
	ek := getEventKey(ctx, eventID)

	keys, err := datastore.NewQuery(eventSnapshotKind).Ancestor(ek).Order("Timestamp").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	es := make([]*EventSnapshot, len(keys))
	err = nds.GetMulti(ctx, keys, es)
	if err != nil {
		return nil, err
	}

	return es, nil
}
