package apicache

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
)

const (
	keyPrefix = "TDAPI:"

	renderEventKeyPrefix  = keyPrefix + "RE1:"
	renderEventExpiryTime = 3 * time.Hour
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
	_ = memcache.Add(ctx, &memcache.Item{
		Key:        renderEventKey(eid),
		Value:      data,
		Expiration: renderEventExpiryTime,
	})
}
