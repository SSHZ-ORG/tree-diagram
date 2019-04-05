package apicache

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

const keyPrefix = "TDAPI:"

type keyFunc func(string) string

func getInternal(ctx context.Context, id string, f keyFunc) []byte {
	if item, err := memcache.Get(ctx, f(id)); err == nil {
		return item.Value
	}
	return nil
}

func putInternal(ctx context.Context, id string, data []byte, f keyFunc) {
	_ = memcache.Set(ctx, &memcache.Item{
		Key:   f(id),
		Value: data,
	})
}

// Errors wrapped.
func clearInternal(ctx context.Context, ids []string, f keyFunc) error {
	var keys []string
	for _, id := range ids {
		keys = append(keys, f(id))
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
