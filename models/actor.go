package models

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Actor struct {
	ID   string `datastore:",noindex"` // Use Key
	Name string

	LastFavoriteCount int

	// The last time we see FavoriteCount changed. Does not care about Name change. Ignored in Equal().
	LastUpdateTime time.Time
}

func (a Actor) debugName() string {
	return fmt.Sprintf("[%s] %s", a.ID, a.Name)
}

func (e *Actor) Equal(o *Actor) bool {
	if e != nil && o != nil {
		if e.ID != o.ID || e.Name != o.Name {
			return false
		}
		if e.LastFavoriteCount != o.LastFavoriteCount {
			return false
		}
		return true
	}
	return e == o
}

const actorKind = "Actor"

func getActorKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, actorKind, id, 0, nil)
}

// Errors wrapped.
func EnsureActors(ctx context.Context, actors map[string]string) (map[string]*datastore.Key, error) {
	var keys []*datastore.Key
	var as []*Actor
	keysMap := make(map[string]*datastore.Key)
	for id, name := range actors {
		key := getActorKey(ctx, id)
		keys = append(keys, key)
		as = append(as, &Actor{ID: id, Name: name})
		keysMap[id] = key
	}

	unused := make([]*Actor, len(keys))
	err := nds.GetMulti(ctx, keys, unused)
	if err == nil {
		// All actors are already there. Just return them back.
		return keysMap, nil
	}

	// We got back some errors, put the non-existing ones.
	var keysToPut []*datastore.Key
	var asToPut []*Actor
	if me, ok := err.(appengine.MultiError); ok {
		for i, e := range me {
			if e == nil {
				// This item is already there.
				continue
			}
			if e != datastore.ErrNoSuchEntity {
				// Something else happened. Rethrow it.
				return nil, errors.Wrap(err, "nds.GetMulti returned error other than NoSuchEntity")
			}
			keysToPut = append(keysToPut, keys[i])
			asToPut = append(asToPut, as[i])
		}

		_, err := nds.PutMulti(ctx, keysToPut, asToPut)
		return keysMap, errors.Wrap(err, "nds.PutMulti failed")
	}

	// WTF?
	return nil, errors.Wrap(err, "nds.GetMulti returned error that is not a MultiError")
}

// Errors wrapped.
func getActors(ctx context.Context, keys []*datastore.Key) ([]*Actor, error) {
	as := make([]*Actor, len(keys))
	err := nds.GetMulti(ctx, keys, as)
	return as, errors.Wrap(err, "nds.GetMulti failed")
}

// Get actors of the given IDs. Not found ones will be ignored.
// Errors wrapped.
func GetActorMap(ctx context.Context, actorIDs []string) (map[string]*Actor, error) {
	var keys []*datastore.Key
	for _, id := range actorIDs {
		keys = append(keys, getActorKey(ctx, id))
	}

	as := make([]*Actor, len(keys))
	err := nds.GetMulti(ctx, keys, as)

	actorMap := make(map[string]*Actor)

	if err == nil {
		// We know all of them.
		for _, a := range as {
			actorMap[a.ID] = a
		}
		return actorMap, nil
	}

	// We got back some errors, just return the existing ones.
	if me, ok := err.(appengine.MultiError); ok {
		for i, e := range me {
			if e != nil {
				if e == datastore.ErrNoSuchEntity {
					// Just skip.
					continue
				} else {
					// Something else happened. Rethrow it.
					return nil, errors.Wrap(err, "nds.GetMulti returned error other than NoSuchEntity")
				}
			}

			// Good, we know this actor.
			actorMap[as[i].ID] = as[i]
		}
		return actorMap, nil
	}

	// WTF?
	return nil, errors.Wrap(err, "nds.GetMulti returned error that is not a MultiError")
}

func UpdateActors(ctx context.Context, oas, nas []*Actor) error {
	var (
		actorKeys    []*datastore.Key
		snapshotKeys []*datastore.Key
		snapshots    []*ActorSnapshot
	)
	for i, na := range nas {
		ak := getActorKey(ctx, na.ID)
		actorKeys = append(actorKeys, ak)
		if sk, snapshot := maybeCreateActorSnapshot(ctx, ak, oas[i], na); sk != nil {
			snapshotKeys = append(snapshotKeys, sk)
			snapshots = append(snapshots, snapshot)
		}
	}

	// Before we can use transaction, put snapshots first.
	// Because we would get into bad state if actors were put successful, but snapshots failed.
	if _, err := nds.PutMulti(ctx, snapshotKeys, snapshots); err != nil {
		return errors.Wrap(err, "nds.PutMulti failed")
	}
	if _, err := nds.PutMulti(ctx, actorKeys, nas); err != nil {
		return errors.Wrap(err, "nds.PutMulti failed")
	}
	return nil
}

type RenderActorResponse struct {
	KnownEventCount int                      `json:"knownEventCount"`
	Snapshots       []*FrontendActorSnapshot `json:"snapshots"`
}

// Errors wrapped.
func PrepareRenderActorResponse(ctx context.Context, actorID string) (*RenderActorResponse, error) {
	key := getActorKey(ctx, actorID)

	response := &RenderActorResponse{
		Snapshots: make([]*FrontendActorSnapshot, 0),
	}

	kec, err := datastore.NewQuery(eventKind).KeysOnly().Filter("Actors =", key).Count(ctx)
	if err != nil {
		return response, errors.Wrap(err, "Counting events failed")
	}
	response.KnownEventCount = kec

	snapshots, err := getFrontendSnapshotsForActor(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(snapshots) > 0 {
		response.Snapshots = snapshots
	}

	return response, nil
}
