// Package appctx defines a function that securely sets and gets
// the request scope data required by the application.
package appctx

import (
	"context"
)

type contextKey string

const requestIDKey contextKey = "request-id"

// WithRequestID returns a copy of the parent context with the requestID set.
func WithRequestID(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, requestIDKey, requestID)
}

// RequestID returns requestID from the context.
// If it does not exists, returns an empty string.
func RequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}
