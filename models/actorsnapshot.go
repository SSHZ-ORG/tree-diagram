package models

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type ActorSnapshot struct {
	ActorID       string `datastore:",noindex"` // Use Key.Parent
	Timestamp     time.Time
	FavoriteCount int `datastore:",noindex"` // Don't query Snapshots directly.
}

const actorSnapshotKind = "ActorSnapshot"

// This returns nil, nil if we should not create a snapshot.
func maybeCreateActorSnapshot(ctx context.Context, ak *datastore.Key, oa, na *Actor) (*datastore.Key, *ActorSnapshot) {
	// For actors, we won't take snapshot if they don't exist yet - oa will not be nil.
	if oa.LastFavoriteCount == na.LastFavoriteCount && !oa.LastUpdateTime.IsZero() {
		return nil, nil
	}

	s := &ActorSnapshot{
		ActorID:       na.ID,
		Timestamp:     na.LastUpdateTime,
		FavoriteCount: na.LastFavoriteCount,
	}

	log.Debugf(ctx, "Taking snapshot for actor %s (%d -> %d)", na.debugName(), oa.LastFavoriteCount, na.LastFavoriteCount)
	return datastore.NewIncompleteKey(ctx, actorSnapshotKind, ak), s
}
