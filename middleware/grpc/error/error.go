// Package grpc_error is a generic server-side gRPC middleware
// that converts standard error to gRPC error.
package grpc_error

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/takuoki/golib/apperr"
	"github.com/takuoki/golib/applog"
)

// UnaryServerInterceptor returns a gRPC middleware that converts standard error to gRPC error.
func UnaryServerInterceptor(domain, internalServerErrorCode string, logger applog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			e, ok := apperr.Extract(err)
			if !ok {
				e = newInternalServerError(internalServerErrorCode, err)
			}
			if e.Log() != "" {
				logger.Error(ctx, e.Log())
			}

			return nil, e.GRPCError(domain)
		}
		return resp, nil
	}
}

func newInternalServerError(code string, err error) apperr.Err {
	return apperr.NewServerError(
		codes.Internal,
		code,
		"internal server error",
		err.Error(),
	)
}
