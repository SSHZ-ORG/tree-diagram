package models

import (
	"time"

	"github.com/pkg/errors"
	"github.com/qedus/nds"
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
func CompressSnapshots(ctx context.Context, eventID string) error {
	eventKey := getEventKey(ctx, eventID)
	log.Debugf(ctx, "Starting compress snapshots for event %s", eventID)

	snapshotKeys, snapshots, err := getNonCompressedSnapshotsForEvent(ctx, eventKey)
	if err != nil {
		return err
	}
	log.Debugf(ctx, "Got %d uncompressed snapshots", len(snapshotKeys))
	if len(snapshotKeys) == 0 {
		// Nothing to compress.
		return nil
	}

	err = nds.RunInTransaction(ctx, func(ctx context.Context) error {
		latestCSKey, latestCS, err := getLatestCompressedSnapshot(ctx, eventKey)
		if err != nil {
			return err
		}
		log.Debugf(ctx, "Last compressed snapshot: %v", latestCSKey)

		entityCountTag := ""

		var keysToPut []*datastore.Key
		var csToPut []*compressedEventSnapshot

		if latestCS.shouldCompress(snapshots[0]) {
			// We should compress the first new snapshot into the last compressed snapshot.
			latestCS.compress(snapshots[0])
			keysToPut = append(keysToPut, latestCSKey)
			csToPut = append(csToPut, latestCS)
			entityCountTag = " (-1)"
		} else {
			// Just create a new compressed snapshot.
			keysToPut = append(keysToPut, datastore.NewIncompleteKey(ctx, compressedEventSnapshotKind, eventKey))
			csToPut = append(csToPut, toCompressedEventSnapshot(snapshots[0]))
		}

		// We now have processed the first snapshot to compress, and have at least one cs in the slice.
		for _, s := range snapshots[1:] {
			if csToPut[len(csToPut)-1].shouldCompress(s) {
				csToPut[len(csToPut)-1].compress(s)
			} else {
				keysToPut = append(keysToPut, datastore.NewIncompleteKey(ctx, compressedEventSnapshotKind, eventKey))
				csToPut = append(csToPut, toCompressedEventSnapshot(s))
			}
		}

		log.Debugf(ctx, "Compressed to %d%s entities", len(keysToPut), entityCountTag)
		if err := nds.DeleteMulti(ctx, snapshotKeys); err != nil {
			return errors.Wrap(err, "nds.DeleteMulti failed")
		}
		if _, err := nds.PutMulti(ctx, keysToPut, csToPut); err != nil {
			return errors.Wrap(err, "nds.PutMulti failed")
		}
		return nil
	}, nil)

	return errors.Wrap(err, "nds.RunInTransaction failed")
}

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

func (c *compressedEventSnapshot) shouldCompress(s *EventSnapshot) bool {
	if c == nil {
		return false
	}
	if c.EventID != s.EventID {
		return false
	}
	if c.NoteCount != s.NoteCount {
		return false
	}
	if len(s.Actors) > 0 {
		return false
	}
	return true
}

func (c *compressedEventSnapshot) compress(s *EventSnapshot) {
	c.Timestamps = append(c.Timestamps, s.Timestamp)
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

func toCompressedEventSnapshot(s *EventSnapshot) *compressedEventSnapshot {
	return &compressedEventSnapshot{
		EventID:    s.EventID,
		Timestamps: []time.Time{s.Timestamp},
		NoteCount:  s.NoteCount,
		Actors:     s.Actors,
	}
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

func newCESFromEvent(ctx context.Context, e *Event, eventKey *datastore.Key) (*datastore.Key, *compressedEventSnapshot) {
	return datastore.NewIncompleteKey(ctx, compressedEventSnapshotKind, eventKey), &compressedEventSnapshot{
		EventID:    e.ID,
		Timestamps: []time.Time{e.LastUpdateTime},
		NoteCount:  e.LastNoteCount,
		Actors:     e.Actors,
	}
}

// Get key and CES which should then be Put into the datastore.
// Errors wrapped.
func createOrUpdateCES(ctx context.Context, oe, ne *Event, eventKey *datastore.Key) (*datastore.Key, *compressedEventSnapshot, error) {
	if shouldCreateNewCES(oe, ne) {
		key, ces := newCESFromEvent(ctx, ne, eventKey)

		log.Debugf(ctx, "Creating new CES for event %s (%d -> %d)", ne.debugName(), oe.LastNoteCount, ne.LastNoteCount)
		return key, ces, nil
	} else {
		cesKey, ces, err := getLatestCompressedSnapshot(ctx, eventKey)
		if err != nil {
			return nil, nil, err
		}
		ces.Timestamps = append(ces.Timestamps, ne.LastUpdateTime)

		log.Debugf(ctx, "Appending to CES %+v for event %s (%d -> %d)", cesKey, ne.debugName(), oe.LastNoteCount, ne.LastNoteCount)
		return cesKey, ces, nil
	}
}
