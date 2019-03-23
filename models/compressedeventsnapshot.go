package models

import (
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type compressedEventSnapshot struct {
	EventID    string `datastore:",noindex"` // Use Key.Parent
	Timestamps []time.Time
	NoteCount  int              `datastore:",noindex"`
	Actors     []*datastore.Key `datastore:",noindex"`
}

const compressedEventSnapshotKind = "CompressedEventSnapshot"

// Errors wrapped.
// This bypasses memcache.
func getCompressedSnapshots(ctx context.Context, eventKey *datastore.Key) ([]*compressedEventSnapshot, error) {
	var css []*compressedEventSnapshot
	_, err := datastore.NewQuery(compressedEventSnapshotKind).Ancestor(eventKey).Order("Timestamps").GetAll(ctx, &css)
	if err != nil {
		return nil, errors.Wrap(err, "datastore query failed")
	}

	return css, nil
}

// Errors wrapped.
// Returns nil if there is no compressed snapshot for the event yet.
// This bypasses memcache.
func getLatestCompressedSnapshot(ctx context.Context, eventKey *datastore.Key) (*datastore.Key, *compressedEventSnapshot, error) {
	var css []*compressedEventSnapshot
	keys, err := datastore.NewQuery(compressedEventSnapshotKind).Ancestor(eventKey).Order("-Timestamps").Limit(1).GetAll(ctx, &css)
	if err != nil {
		return nil, nil, errors.Wrap(err, "datastore query failed")
	}

	if len(keys) == 0 {
		return nil, nil, nil
	}
	return keys[0], css[0], nil
}

func (c *compressedEventSnapshot) decompress() []*EventSnapshot {
	var ss []*EventSnapshot
	for _, ts := range c.Timestamps {
		ss = append(ss, &EventSnapshot{
			EventID:   c.EventID,
			Timestamp: ts,
			NoteCount: c.NoteCount,
		})
	}
	ss[0].Actors = c.Actors
	return ss
}

func (c *compressedEventSnapshot) isConsistent(e *Event) bool {
	if c == nil || e == nil {
		return false
	}
	if c.NoteCount != e.LastNoteCount {
		return false
	}
	if len(c.Actors) > 0 && !areKeysSetsEqual(c.Actors, e.Actors) {
		return false
	}
	return true
}

func shouldCreateNewCES(oe, ne *Event) bool {
	if oe == nil || ne == nil {
		return true
	}
	if oe.LastNoteCount != ne.LastNoteCount {
		return true
	}
	if !areKeysSetsEqual(oe.Actors, ne.Actors) {
		return true
	}
	return false
}

func newCESFromEvent(ctx context.Context, oe, ne *Event, eventKey *datastore.Key) (*datastore.Key, *compressedEventSnapshot) {
	ces := &compressedEventSnapshot{
		EventID:    ne.ID,
		Timestamps: []time.Time{ne.LastUpdateTime},
		NoteCount:  ne.LastNoteCount,
	}
	if oe == nil || !areKeysSetsEqual(oe.Actors, ne.Actors) {
		ces.Actors = ne.Actors
	}
	return datastore.NewIncompleteKey(ctx, compressedEventSnapshotKind, eventKey), ces
}

// Get key and CES which should then be Put into the datastore.
// Errors wrapped.
func createOrUpdateCES(ctx context.Context, oe, ne *Event, eventKey *datastore.Key) (*datastore.Key, *compressedEventSnapshot, error) {
	var key *datastore.Key
	var ces *compressedEventSnapshot

	if !shouldCreateNewCES(oe, ne) {
		var err error
		key, ces, err = getLatestCompressedSnapshot(ctx, eventKey)
		if err != nil {
			return nil, nil, err
		}

		if ces.isConsistent(ne) {
			ces.Timestamps = append(ces.Timestamps, ne.LastUpdateTime)
			log.Debugf(ctx, "Appending to CES %+v for event %s (%d -> %d)", key, ne.debugName(), oe.LastNoteCount, ne.LastNoteCount)
		} else {
			key, ces = nil, nil
			log.Criticalf(ctx, "Inconsistent CES %+v detected!", key)
		}
	}

	if key == nil {
		key, ces = newCESFromEvent(ctx, oe, ne, eventKey)

		lastNoteCount := 0
		if oe != nil {
			lastNoteCount = oe.LastNoteCount
		}
		log.Debugf(ctx, "Creating new CES for event %s (%d -> %d)", ne.debugName(), lastNoteCount, ne.LastNoteCount)
	}
	return key, ces, nil
}
