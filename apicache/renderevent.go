package apicache

import "golang.org/x/net/context"

const renderEventKeyPrefix = keyPrefix + "RE4:"

func renderEventKey(eid string) string {
	return renderEventKeyPrefix + eid
}

func GetRenderEvent(ctx context.Context, eid string) []byte {
	return getInternal(ctx, []string{eid}, renderEventKey)[eid]
}

func PutRenderEvent(ctx context.Context, eid string, data []byte) {
	putInternal(ctx, map[string][]byte{eid: data}, renderEventKey)
}

// Errors wrapped.
func ClearRenderEvents(ctx context.Context, eids []string) error {
	return clearInternal(ctx, eids, renderEventKey)
}
