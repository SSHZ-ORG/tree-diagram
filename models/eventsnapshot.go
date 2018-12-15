package models

import (
	"time"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type EventSnapshot struct {
	Timestamp time.Time

	NoteCount int
	Actors    []*datastore.Key
}

const eventSnapshotKind = "EventSnapshot"

// Is a snapshot needed?
func needSnapshot(old, new *Event) bool {
	if old.LastNoteCount != new.LastNoteCount {
		return true
	}

	return !areKeysSetsEqual(old.Actors, new.Actors)
}

// Take a snapshot of the event. Can be called in a transaction.
func takeSnapshot(ctx context.Context, ek *datastore.Key, e *Event) error {
	key := datastore.NewIncompleteKey(ctx, eventSnapshotKind, ek)
	s := &EventSnapshot{
		Timestamp: e.LastUpdateTime,
		NoteCount: e.LastNoteCount,
		Actors:    e.Actors,
	}

	_, err := nds.Put(ctx, key, s)
	return err
}
