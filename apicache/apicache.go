package apicache

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

const keyPrefix = "TDAPI:"

type keyFunc func(string) string

func getInternal(ctx context.Context, ids []string, f keyFunc) map[string][]byte {
	var keys []string
	for _, id := range ids {
		keys = append(keys, f(id))
	}
	items, _ := memcache.GetMulti(ctx, keys)

	m := make(map[string][]byte)
	for _, id := range ids {
		if item, ok := items[f(id)]; ok {
			m[id] = item.Value
		}
	}
	return m
}

func putInternal(ctx context.Context, data map[string][]byte, f keyFunc) {
	var items []*memcache.Item
	for id, v := range data {
		items = append(items, &memcache.Item{
			Key:   f(id),
			Value: v,
		})
	}
	_ = memcache.SetMulti(ctx, items)
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
