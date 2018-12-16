package models

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type EventSnapshot struct {
	Timestamp time.Time

	// Will only snapshot changed fields.
	NoteCount int
	Actors    []*datastore.Key
}

const eventSnapshotKind = "EventSnapshot"

func maybeCreateSnapshot(ctx context.Context, ek *datastore.Key, oe, ne *Event) (*datastore.Key, *EventSnapshot) {
	shouldTake := false
	s := &EventSnapshot{Timestamp: ne.LastUpdateTime}

	if oe == nil {
		oe = &Event{}
	}

	if oe.LastNoteCount != ne.LastNoteCount {
		s.NoteCount = ne.LastNoteCount
		shouldTake = true
	}

	if !areKeysSetsEqual(oe.Actors, ne.Actors) {
		s.Actors = ne.Actors
		shouldTake = true
	}

	if shouldTake {
		log.Debugf(ctx, "Taking snapshot for event %s.", ne.debugName())
		return datastore.NewIncompleteKey(ctx, eventSnapshotKind, ek), s
	} else {
		log.Debugf(ctx, "Snapshot skipped for event %s.", ne.debugName())
		return nil, nil
	}
}
