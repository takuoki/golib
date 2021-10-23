// Package appctx defines a function that securely sets and gets
// the request scope data required by the application.
package appctx

import (
	"context"
)

type contextKey string

// List of contextKey
const (
	requestIDKey     contextKey = "request-id"
	userIDKey        contextKey = "user-id"
	authorizationKey contextKey = "authorization"
)

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

// WithUserID returns a copy of the parent context with the userID set.
func WithUserID(parent context.Context, userID string) context.Context {
	return context.WithValue(parent, userIDKey, userID)
}

// UserID returns userID from the context.
// If it does not exists, returns an empty string.
func UserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// WithAuthorization returns a copy of the parent context with the authorization set.
func WithAuthorization(parent context.Context, authorization string) context.Context {
	return context.WithValue(parent, authorizationKey, authorization)
}

// Authorization returns authorization from the context.
// If it does not exists, returns an empty string.
func Authorization(ctx context.Context) string {
	if authorization, ok := ctx.Value(authorizationKey).(string); ok {
		return authorization
	}
	return ""
}
