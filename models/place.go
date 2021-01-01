package models

import (
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/protobuf/proto"
)

type Place struct {
	ID   string `datastore:",noindex"` // Use Key
	Name string

	ShouldIgnore    bool
	DefaultCapacity int // GAE doesn't care. All int* types are 64-bit.
}

const (
	placeKind = "Place"
)

func getPlaceKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, placeKind, id, 0, nil)
}

// Errors wrapped.
func EnsurePlaces(ctx context.Context, places map[string]string) (map[string]*datastore.Key, error) {
	var keys []*datastore.Key
	var ps []*Place
	keysMap := make(map[string]*datastore.Key)
	for id, name := range places {
		key := getPlaceKey(ctx, id)
		keys = append(keys, key)
		ps = append(ps, &Place{ID: id, Name: name})
		keysMap[id] = key
	}

	unused := make([]*Place, len(keys))
	err := nds.GetMulti(ctx, keys, unused)
	if err == nil {
		// All places are already there. Just return them back.
		return keysMap, nil
	}

	// We got back some errors, put the non-existing ones.
	var keysToPut []*datastore.Key
	var psToPut []*Place
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
			psToPut = append(psToPut, ps[i])
		}

		_, err := nds.PutMulti(ctx, keysToPut, psToPut)
		return keysMap, errors.Wrap(err, "nds.PutMulti failed")
	}

	// WTF?
	return nil, errors.Wrap(err, "nds.GetMulti returned error that is not a MultiError")
}

type RenderPlaceResponse struct {
	KnownEventCount int `json:"knownEventCount"`
}

// Errors wrapped.
func PrepareRenderPlaceResponse(ctx context.Context, placeID string) (*pb.RenderPlaceResponse, error) {
	key := getPlaceKey(ctx, placeID)

	response := &pb.RenderPlaceResponse{}

	kec, err := datastore.NewQuery(eventKind).KeysOnly().Filter("Place =", key).Count(ctx)
	if err != nil {
		return response, errors.Wrap(err, "Counting events failed")
	}
	response.KnownEventCount = proto.Int32(int32(kec))

	return response, nil
}
