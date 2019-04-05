package apicache

import "golang.org/x/net/context"

const renderEventKeyPrefix = keyPrefix + "RE2:"

func renderEventKey(eid string) string {
	return renderEventKeyPrefix + eid
}

func GetRenderEvent(ctx context.Context, eid string) []byte {
	return getInternal(ctx, eid, renderEventKey)
}

func PutRenderEvent(ctx context.Context, eid string, data []byte) {
	putInternal(ctx, eid, data, renderEventKey)
}

// Errors wrapped.
func ClearRenderEvents(ctx context.Context, eids []string) error {
	return clearInternal(ctx, eids, renderEventKey)
}
