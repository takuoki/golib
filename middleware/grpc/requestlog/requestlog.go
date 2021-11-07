// Package grpc_requestlog is a generic server-side gRPC middleware
// that outputs request log.
// If the request ID is not specified, it will be automatically generated.
package grpc_requestlog

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/applog"
)

const (
	contentTypeKey = "content-type"
	userAgentKey   = "user-agent"
)

// UnaryServerInterceptor returns a gRPC middleware that outputs request logs.
func UnaryServerInterceptor(logger applog.Logger, opt ...Option) grpc.UnaryServerInterceptor {

	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// Set request ID to context.
		reqID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			r := md.Get(opts.requestIDKey)
			if len(r) > 0 {
				reqID = r[0]
			}
		}
		if reqID == "" {
			id, err := opts.requestIDFunc()
			if err != nil {
				logger.Warnf(ctx, "failed to create new request ID: %v", err)
			}
			reqID = id
		}
		ctx = appctx.WithRequestID(ctx, reqID)

		// Create label.
		label := map[string]string{
			"service_method": info.FullMethod,
		}
		if pr, ok := peer.FromContext(ctx); ok {
			label["ip_address"] = pr.Addr.String()
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ct := md.Get(contentTypeKey)
			if len(ct) > 0 {
				label["content_type"] = ct[0]
			}
			ua := md.Get(userAgentKey)
			if len(ua) > 0 {
				label["user_agent"] = ua[0]
			}
		}

		logger.Print(ctx, applog.InfoLevel, "request log", label)

		return handler(ctx, req)
	}
}
