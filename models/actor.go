package models

import (
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Actor struct {
	ID   string
	Name string
}

const actorKind = "Actor"

func EnsureActors(ctx context.Context, actors map[string]string) (map[string]*datastore.Key, error) {
	var keys []*datastore.Key
	var as []*Actor
	keysMap := make(map[string]*datastore.Key)
	for id, name := range actors {
		key := datastore.NewKey(ctx, actorKind, id, 0, nil)
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
				return nil, err
			}
			keysToPut = append(keysToPut, keys[i])
			asToPut = append(asToPut, as[i])
		}

		_, err := nds.PutMulti(ctx, keysToPut, asToPut)
		return keysMap, err
	}

	// WTF?
	return nil, err
}
