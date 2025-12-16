package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func AuthMetadata(ctx context.Context, r *http.Request) metadata.MD {
	if uid, ok := ctx.Value(userIDKey{}).(string); ok && uid != "" {
		return metadata.Pairs("x-user-id", uid)
	}
	return nil
}
