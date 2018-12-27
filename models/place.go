package models

import (
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Place struct {
	ID   string `datastore:",noindex"` // Use Key
	Name string

	ShouldIgnore    bool
	DefaultCapacity int // GAE doesn't care. All int* types are 64-bit.
}

const (
	placeKind = "Place"

	renderPlaceMaxTopEvents = 20
)

func getPlaceKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, placeKind, id, 0, nil)
}

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
				return nil, err
			}
			keysToPut = append(keysToPut, keys[i])
			psToPut = append(psToPut, ps[i])
		}

		_, err := nds.PutMulti(ctx, keysToPut, psToPut)
		return keysMap, err
	}

	// WTF?
	return nil, err
}

type RenderPlaceResponse struct {
	TopEvents []FrontendEvent `json:"top_events"`
}

func PrepareRenderPlaceResponse(ctx context.Context, placeID string) (*RenderPlaceResponse, error) {
	key := getPlaceKey(ctx, placeID)

	response := &RenderPlaceResponse{
		TopEvents: make([]FrontendEvent, 0), // So json does not make it null.
	}

	keys, err := datastore.NewQuery(eventKind).KeysOnly().Filter("Place =", key).Order("-LastNoteCount").Limit(renderPlaceMaxTopEvents).GetAll(ctx, nil)
	if err != nil {
		return response, err
	}

	es := make([]*Event, len(keys))
	err = nds.GetMulti(ctx, keys, es)
	if err != nil {
		return response, err
	}

	for _, e := range es {
		response.TopEvents = append(response.TopEvents, e.ToFrontendEvent())
	}

	return response, nil
}
