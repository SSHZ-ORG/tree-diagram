package models

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/protobuf/proto"
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
// oa cannot be nil.
func maybeCreateActorSnapshot(ctx context.Context, ak *datastore.Key, oa, na *Actor) (*datastore.Key, *ActorSnapshot) {
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

// Errors wrapped.
func getFrontendSnapshotsForActor(ctx context.Context, ak *datastore.Key) ([]*pb.RenderActorsResponse_ResponseItem_Snapshot, error) {
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

	var fass []*pb.RenderActorsResponse_ResponseItem_Snapshot
	for _, s := range snapshots {
		fass = append(fass, &pb.RenderActorsResponse_ResponseItem_Snapshot{
			Date:          utils.ToProtoDate(civil.DateOf(s.Timestamp.In(utils.JST()))),
			FavoriteCount: proto.Int32(int32(s.FavoriteCount)),
		})
	}
	return fass, nil
}

func OneoffBackfillModelVersion(ctx context.Context, cursor string) (string, error) {
	q := datastore.NewQuery(actorSnapshotKind).KeysOnly()

	if cursor != "" {
		c, err := datastore.DecodeCursor(cursor)
		if err != nil {
			panic(err)
		}
		q = q.Start(c)
	}
	q = q.Limit(25)

	var keys []*datastore.Key
	it := q.Run(ctx)
	key, err := it.Next(nil)
	for err == nil {
		keys = append(keys, key)
		key, err = it.Next(nil)
	}
	if err != datastore.Done {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	err = nds.RunInTransaction(ctx, func(ctx context.Context) error {
		ass := make([]*ActorSnapshot, len(keys))
		if err := nds.GetMulti(ctx, keys, ass); err != nil {
			return errors.Wrap(err, "nds.GetMulti failed")
		}

		var keysToPut []*datastore.Key
		var assToPut []*ActorSnapshot

		for i, as := range ass {
			if as.ModelVersion == 0 {
				as.ModelVersion = actorSnapshotCurrentModelVersion
				keysToPut = append(keysToPut, keys[i])
				assToPut = append(assToPut, as)
			}
		}

		log.Debugf(ctx, "Updating %d ASs", len(keysToPut))
		_, err := nds.PutMulti(ctx, keysToPut, assToPut)
		return errors.Wrap(err, "nds.PutMulti failed")
	}, &datastore.TransactionOptions{XG: true})

	if err != nil {
		return "", err
	}
	newCursor, err := it.Cursor()
	if err != nil {
		panic(err)
	}
	return newCursor.String(), nil
}
