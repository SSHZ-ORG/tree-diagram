package models

import (
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Actor struct {
	ID   string `datastore:",noindex"` // Use Key
	Name string
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

// Error wrapped
func GetActors(ctx context.Context, keys []*datastore.Key) ([]*Actor, error) {
	as := make([]*Actor, len(keys))
	err := nds.GetMulti(ctx, keys, as)
	return as, errors.Wrap(err, "nds.GetMulti failed")
}

type RenderActorResponse struct {
	KnownEventCount int `json:"knownEventCount"`
}

// Errors wrapped.
func PrepareRenderActorResponse(ctx context.Context, actorID string) (*RenderActorResponse, error) {
	key := getActorKey(ctx, actorID)

	response := &RenderActorResponse{}

	kec, err := datastore.NewQuery(eventKind).KeysOnly().Filter("Actors =", key).Count(ctx)
	if err != nil {
		return response, errors.Wrap(err, "Counting events failed")
	}
	response.KnownEventCount = kec

	return response, nil
}
