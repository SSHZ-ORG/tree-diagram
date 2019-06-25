package apicache

import "golang.org/x/net/context"

const renderActorKeyPrefix = keyPrefix + "RA1:"

func renderActorKey(aid string) string {
	return renderActorKeyPrefix + aid
}

func GetRenderActor(ctx context.Context, aid string) []byte {
	return getInternal(ctx, []string{aid}, renderActorKey)[aid]
}

func PutRenderActor(ctx context.Context, aid string, data []byte) {
	putInternal(ctx, map[string][]byte{aid: data}, renderActorKey)
}

// Errors wrapped.
func ClearRenderActors(ctx context.Context, aids []string) error {
	return clearInternal(ctx, aids, renderActorKey)
}
