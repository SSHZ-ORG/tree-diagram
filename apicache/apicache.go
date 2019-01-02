package apicache

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

const (
	keyPrefix = "TDAPI:"

	renderEventKeyPrefix = keyPrefix + "RE1:"
)

func renderEventKey(eid string) string {
	return renderEventKeyPrefix + eid
}

func GetRenderEvent(ctx context.Context, eid string) []byte {
	if item, err := memcache.Get(ctx, renderEventKey(eid)); err == nil {
		return item.Value
	}
	return nil
}

func PutRenderEvent(ctx context.Context, eid string, data []byte) {
	_ = memcache.Set(ctx, &memcache.Item{
		Key:   renderEventKey(eid),
		Value: data,
	})
}

// Errors wrapped.
func ClearRenderEvents(ctx context.Context, eids []string) error {
	var keys []string
	for _, eid := range eids {
		keys = append(keys, renderEventKey(eid))
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
