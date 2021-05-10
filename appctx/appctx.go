package appctx

import (
	"context"
)

type contextKey string

const requestIDContextKey contextKey = "request-id"

// SetRequestID sets requestID in the context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

// GetRequestID returns requestID from the context.
// If it does not exists, returns an empty string.
func GetRequestID(ctx context.Context) string {
	v := ctx.Value(requestIDContextKey)
	requestID, ok := v.(string)
	if !ok {
		return ""
	}
	return requestID
}
