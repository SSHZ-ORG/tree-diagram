package models

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type EventSnapshot struct {
	Timestamp time.Time

	// Recorded even if not changed.
	NoteCount int

	// Recorded only if changed.
	Actors []*datastore.Key
}

const eventSnapshotKind = "EventSnapshot"

func maybeCreateSnapshot(ctx context.Context, ek *datastore.Key, oe, ne *Event) (*datastore.Key, *EventSnapshot) {
	s := &EventSnapshot{
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
