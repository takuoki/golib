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
	if requestID, ok := ctx.Value(requestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}
