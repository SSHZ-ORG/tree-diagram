package models

import (
	"fmt"
	"time"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
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

func (e Event) debugName() string {
	return fmt.Sprintf("[%s] %s", e.ID, e.Name)
}

const eventKind = "Event"

// Insert or update events. This automatically takes snapshots if needed.
func InsertOrUpdateEvents(ctx context.Context, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	var keys []*datastore.Key
	for _, e := range events {
		keys = append(keys, datastore.NewKey(ctx, eventKind, e.ID, 0, nil))
	}

	oes := make([]*Event, len(events))
	err := nds.GetMulti(ctx, keys, oes)
	if err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for _, e := range me {
				if e != nil && e != datastore.ErrNoSuchEntity {
					// Something else happened. Rethrow it.
					return err
				}
			}
		} else {
			return err
		}
	}

	var snapshotKeys []*datastore.Key
	var snapshots []*EventSnapshot
	for i, e := range events {
		sk, s := maybeCreateSnapshot(ctx, keys[i], oes[i], e)
		if sk != nil {
			snapshotKeys = append(snapshotKeys, sk)
			snapshots = append(snapshots, s)
		}
	}

	// We always update events even if no change, to update the timestamp.
	if _, err := nds.PutMulti(ctx, keys, events); err != nil {
		return err
	}
	if _, err := nds.PutMulti(ctx, snapshotKeys, snapshots); err != nil {
		return err
	}

	return nil
}
