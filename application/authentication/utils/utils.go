package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type contextKey string

var (
	ContextKeyUID = contextKey("uid")
)

func GetValueFromMetadata(key string, ctx context.Context) string {
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return ""
	}

	values := md[key]

	if len(values) == 0 {
		return ""
	}

	return values[0]
}
