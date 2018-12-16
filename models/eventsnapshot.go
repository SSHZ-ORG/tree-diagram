package models

import (
	"time"

	"github.com/qedus/nds"
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

func maybeTakeSnapshot(ctx context.Context, ek *datastore.Key, oe, ne *Event) error {
	shouldTake := false
	s := &EventSnapshot{Timestamp: ne.LastUpdateTime}

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
		key := datastore.NewIncompleteKey(ctx, eventSnapshotKind, ek)
		_, err := nds.Put(ctx, key, s)
		return err
	} else {
		log.Debugf(ctx, "Snapshot skipped for event %s.", ne.debugName())
		return nil
	}
}
