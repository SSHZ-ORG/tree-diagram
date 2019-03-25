// This is prone to race condition. Do not attempt to multi process the same events at the same time.

package cache

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

const lastCESKeyPrefix = keyPrefix + "LCESK1:"

func lastCESKey(eid string) string {
	return lastCESKeyPrefix + eid
}

// Does not return error. If failed (not found or memcache internal error), just log and return nil.
func GetLastCESKey(ctx context.Context, eventKey *datastore.Key) *datastore.Key {
	if eventKey.Kind() != "Event" { // TODO: don't hardcode this
		log.Criticalf(ctx, "Illegal event Key %+v", eventKey)
		return nil
	}
	eid := eventKey.StringID()
	if item, err := memcache.Get(ctx, lastCESKey(eid)); err == nil {
		k, err := datastore.DecodeKey(string(item.Value))
		if err != nil {
			log.Warningf(ctx, "Invalid last CES Key for event %s from memcache: %+v", eid, err)
		}
		return k
	} else {
		if err != memcache.ErrCacheMiss {
			log.Warningf(ctx, "memcache failed when fetching last CES Key for event %s: %+v", eid, err)
		}
		return nil
	}
}

// Errors wrapped.
func ClearLastCESKeys(ctx context.Context, eventKeys []*datastore.Key) error {
	var keys []string
	for _, k := range eventKeys {
		if k.Kind() != "Event" { // TODO: don't hardcode this
			log.Criticalf(ctx, "Illegal event Key %+v", k)
			continue
		}
		keys = append(keys, lastCESKey(k.StringID()))
	}
	err := memcache.DeleteMulti(ctx, keys)

	if me, ok := err.(appengine.MultiError); ok {
		for _, e := range me {
			if e != nil && e != memcache.ErrCacheMiss {
				// Something else happened. Rethrow it.
				return errors.Wrap(err, "memcache.DeleteMulti returned error other than CacheMiss")
			}
		}
		return nil
	}

	return errors.Wrap(err, "memcache.DeleteMulti returned error that is not a MultiError")
}

// Errors will just be logged.
func SetLastCESKeys(ctx context.Context, keys []*datastore.Key) {
	var items []*memcache.Item
	for _, k := range keys {
		if k.Parent().Kind() != "Event" { // TODO: don't hardcode this
			log.Criticalf(ctx, "Illegal CES Key %+v", k)
			continue
		}
		items = append(items, &memcache.Item{
			Key:   lastCESKey(k.Parent().StringID()),
			Value: []byte(k.Encode()),
		})
	}

	if err := memcache.SetMulti(ctx, items); err != nil {
		log.Warningf(ctx, "memcache failed when setting last CES Keys: %+v", err)
	}
}
