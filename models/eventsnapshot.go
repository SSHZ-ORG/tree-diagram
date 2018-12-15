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

// Take a snapshot of the event. Can be called in a transaction.
func TakeSnapshot(ctx context.Context, ek *datastore.Key, e *Event) error {
	key := datastore.NewIncompleteKey(ctx, eventSnapshotKind, ek)
	s := &EventSnapshot{
		Timestamp: e.LastUpdateTime,
		NoteCount: e.LastNoteCount,
		Actors:    e.Actors,
	}

	_, err := nds.Put(ctx, key, s)
	return err
}
