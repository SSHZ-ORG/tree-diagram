package models

import (
	"time"

	"github.com/SSHZ-ORG/tree-diagram/models/cache"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type compressedEventSnapshot struct {
	EventID      string `datastore:",noindex"` // Use Key.Parent
	Timestamps   []time.Time
	NoteCount    int              `datastore:",noindex"`
	Actors       []*datastore.Key `datastore:",noindex"`
	ModelVersion int
}

const (
	compressedEventSnapshotKind                = "CompressedEventSnapshot"
	compressedEventSnapshotCurrentModelVersion = 2
)

// Errors wrapped.
// This bypasses memcache.
func getCompressedSnapshots(ctx context.Context, eventKey *datastore.Key) ([]*compressedEventSnapshot, error) {
	keys, err := datastore.NewQuery(compressedEventSnapshotKind).Ancestor(eventKey).Order("-Timestamps").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "datastore query failed")
	}

	// We query with ORDER BY Timestamps DESC. Reverse it back.
	// Note that this field is repeated so it works only if all CES time ranges are consecutive.
	for left, right := 0, len(keys)-1; left < right; left, right = left+1, right-1 {
		keys[left], keys[right] = keys[right], keys[left]
	}

	css := make([]*compressedEventSnapshot, len(keys))
	err = nds.GetMulti(ctx, keys, css)
	return css, errors.Wrap(err, "nds.GetMulti failed")
}

// Returns nil if there is no compressed snapshot for the event yet.
// This is prone to race condition. Do not multi process the same events at the same time.
// Errors wrapped.
func getLatestCompressedSnapshotKey(ctx context.Context, eventKey *datastore.Key) (*datastore.Key, error) {
	if maybeCESKey := cache.GetLastCESKey(ctx, eventKey); maybeCESKey != nil {
		return maybeCESKey, nil
	}

	keys, err := datastore.NewQuery(compressedEventSnapshotKind).Ancestor(eventKey).Order("-Timestamps").Limit(1).KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "datastore query failed")
	}

	if len(keys) == 0 {
		return nil, nil
	}
	return keys[0], nil
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
	if len(c.Actors) > 0 && len(e.Actors) > 0 && !areKeysSetsEqual(c.Actors, e.Actors) {
		return false
	}
	return true
}

func (c *compressedEventSnapshot) appendTime(ts time.Time) {
	c.Timestamps = append(c.Timestamps, ts)
	c.ModelVersion = compressedEventSnapshotCurrentModelVersion
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
		EventID:      ne.ID,
		Timestamps:   []time.Time{ne.LastUpdateTime},
		NoteCount:    ne.LastNoteCount,
		ModelVersion: compressedEventSnapshotCurrentModelVersion,
	}
	if oe == nil || !areKeysSetsEqual(oe.Actors, ne.Actors) {
		ces.Actors = ne.Actors
	}
	return datastore.NewIncompleteKey(ctx, compressedEventSnapshotKind, eventKey), ces
}

// If we should append to the last CES, returns its Key. Otherwise returns nil.
// Errors wrapped.
func maybeGetCESKeyToAppend(ctx context.Context, oe, ne *Event, eventKey *datastore.Key) (*datastore.Key, error) {
	if !shouldCreateNewCES(oe, ne) {
		return getLatestCompressedSnapshotKey(ctx, eventKey)
	}
	return nil, nil
}
