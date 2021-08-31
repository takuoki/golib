// Package grpc_error is a generic server-side gRPC middleware
// that converts standard error to gRPC error.
package grpc_error

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takuoki/golib/apperr"
	"github.com/takuoki/golib/applog"
	"github.com/takuoki/golib/notice"
)

// UnaryServerInterceptor returns a gRPC middleware that converts standard error to gRPC error.
func UnaryServerInterceptor(domain, internalServerErrorCode string, logger applog.Logger, notifier notice.Notifier) grpc.UnaryServerInterceptor {
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
			if e.Type() == apperr.ServerError {
				if err := notifier.Error(err); err != nil {
					logger.Errorf(ctx, "failed to send nortification: %v", err)
				}
			}

			st := status.New(codes.Code(e.Status()), e.Message())
			st, _ = st.WithDetails(&errdetails.ErrorInfo{
				Reason: e.Code(),
				Domain: domain,
			})
			return nil, st.Err()
		}
		return resp, nil
	}
}

func newInternalServerError(code string, err error) apperr.Err {
	return apperr.NewServerError(
		int(codes.Internal),
		code,
		"internal server error",
		err.Error(),
	)
}
