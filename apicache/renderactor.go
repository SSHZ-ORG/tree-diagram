package apicache

import "context"

const renderActorKeyPrefix = keyPrefix + "RA2:"

func renderActorKey(aid string) string {
	return renderActorKeyPrefix + aid
}

func GetRenderActor(ctx context.Context, aid []string) map[string][]byte {
	return getInternal(ctx, aid, renderActorKey)
}

func PutRenderActor(ctx context.Context, data map[string][]byte) {
	putInternal(ctx, data, renderActorKey)
}

// Errors wrapped.
func ClearRenderActors(ctx context.Context, aids []string) error {
	return clearInternal(ctx, aids, renderActorKey)
}
