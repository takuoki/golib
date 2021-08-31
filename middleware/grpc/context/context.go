// Package grpc_context is a generic server-side gRPC middleware
// that sets some values to context.
package grpc_context

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/applog"
)

type KeyInfo struct {
	requestIDKey     string
	authorizationKey string
}

type Option func(*KeyInfo)

// RequestIDKey is a key of request ID.
func RequestIDKey(key string) Option {
	return func(k *KeyInfo) {
		k.requestIDKey = key
	}
}

// AuthorizationKey is a key of authorization.
func AuthorizationKey(key string) Option {
	return func(k *KeyInfo) {
		k.authorizationKey = key
	}
}

// UnaryServerInterceptor returns a gRPC middleware that sets some values to context.
func UnaryServerInterceptor(logger applog.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	keyInfo := KeyInfo{}
	for _, opt := range opts {
		opt(&keyInfo)
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// request ID
		reqID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			r := md.Get(keyInfo.requestIDKey)
			if len(r) > 0 {
				reqID = r[0]
			}
		}
		if reqID == "" {
			u, err := uuid.NewRandom()
			if err != nil {
				logger.Warnf(ctx, "failed to create new UUID for request ID in context middleware: %v", err)
			}
			reqID = u.String()
		}
		ctx = appctx.WithRequestID(ctx, reqID)

		// authorization
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			a := md.Get(keyInfo.authorizationKey)
			if len(a) > 0 {
				ctx = appctx.WithAuthorization(ctx, a[0])
			}
		}

		return handler(ctx, req)
	}
}
