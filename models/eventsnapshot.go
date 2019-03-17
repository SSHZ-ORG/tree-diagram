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

func createEventSnapshot(ctx context.Context, ek *datastore.Key, oe, ne *Event, ts time.Time) (*datastore.Key, *EventSnapshot) {
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

	log.Debugf(ctx, "Taking snapshot for event %s (%d -> %d)", ne.debugName(), oe.LastNoteCount, ne.LastNoteCount)
	return datastore.NewIncompleteKey(ctx, eventSnapshotKind, ek), s
}

// Returns all snapshots, including uncompressed ones and decompressed compressed ones.
// The bool returned is whether there are uncompressed snapshots.
// Errors wrapped.
func getSnapshotsForEvent(ctx context.Context, eventKey *datastore.Key) ([]*EventSnapshot, bool, error) {
	cess, err := getCompressedSnapshots(ctx, eventKey)
	if err != nil {
		return nil, false, err
	}

	_, ss, err := getNonCompressedSnapshotsForEvent(ctx, eventKey)
	if err != nil {
		return nil, false, err
	}

	var merged []*EventSnapshot
	for _, ces := range cess {
		merged = append(merged, ces.decompress()...)
	}
	log.Debugf(ctx, "CompressedEventSnapshot: %d/%d (%.2f%%)", len(cess), len(merged), float64(len(cess))/float64(len(merged))*100)

	merged = append(merged, ss...)
	return merged, len(ss) > 0, nil
}

// Errors wrapped.
func getNonCompressedSnapshotsForEvent(ctx context.Context, eventKey *datastore.Key) ([]*datastore.Key, []*EventSnapshot, error) {
	keys, err := datastore.NewQuery(eventSnapshotKind).Ancestor(eventKey).Order("Timestamp").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "datastore query failed")
	}

	es := make([]*EventSnapshot, len(keys))
	err = nds.GetMulti(ctx, keys, es)
	if err != nil {
		return nil, nil, errors.Wrap(err, "nds.GetMulti failed")
	}

	return keys, es, nil
}
