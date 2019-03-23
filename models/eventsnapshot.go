package models

import (
	"time"

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

// Returns all snapshots, decompressed from CES.
// Errors wrapped.
func getSnapshotsForEvent(ctx context.Context, eventKey *datastore.Key) ([]*EventSnapshot, error) {
	cess, err := getCompressedSnapshots(ctx, eventKey)
	if err != nil {
		return nil, err
	}

	var merged []*EventSnapshot
	for _, ces := range cess {
		merged = append(merged, ces.decompress()...)
	}
	log.Debugf(ctx, "CompressedEventSnapshot: %d/%d (%.2f%%)", len(cess), len(merged), float64(len(cess))/float64(len(merged))*100)

	return merged, nil
}
