package models

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type ActorSnapshot struct {
	ActorID       string `datastore:",noindex"` // Use Key.Parent
	Timestamp     time.Time
	FavoriteCount int `datastore:",noindex"` // Don't query Snapshots directly.
	ModelVersion  int
}

const (
	actorSnapshotKind                = "ActorSnapshot"
	actorSnapshotCurrentModelVersion = 1
)

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
		ModelVersion:  actorSnapshotCurrentModelVersion,
	}

	log.Debugf(ctx, "Taking snapshot for actor %s (%d -> %d)", na.debugName(), oa.LastFavoriteCount, na.LastFavoriteCount)
	return datastore.NewIncompleteKey(ctx, actorSnapshotKind, ak), s
}

func (as ActorSnapshot) toFrontendSnapshot() *FrontendActorSnapshot {
	return &FrontendActorSnapshot{
		Date:          civil.DateOf(as.Timestamp.In(utils.JST())),
		FavoriteCount: as.FavoriteCount,
	}
}

type FrontendActorSnapshot struct {
	Date          civil.Date `json:"date"`
	FavoriteCount int        `json:"favoriteCount"`
}

// Errors wrapped.
func getFrontendSnapshotsForActor(ctx context.Context, ak *datastore.Key) ([]*FrontendActorSnapshot, error) {
	keys, err := datastore.NewQuery(actorSnapshotKind).Ancestor(ak).Order("-Timestamp").KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "datastore query failed")
	}

	// We query with ORDER BY Timestamps DESC. Reverse it back.
	for left, right := 0, len(keys)-1; left < right; left, right = left+1, right-1 {
		keys[left], keys[right] = keys[right], keys[left]
	}

	snapshots := make([]*ActorSnapshot, len(keys))
	err = nds.GetMulti(ctx, keys, snapshots)
	if err != nil {
		return nil, errors.Wrap(err, "nds.GetMulti failed")
	}

	var fass []*FrontendActorSnapshot
	for _, s := range snapshots {
		fass = append(fass, s.toFrontendSnapshot())
	}
	return fass, nil
}
