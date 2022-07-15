// Package apperr defines an error that can distinguish
// between a client error and a server error.
// It also defines a function that extracts the defined
// client or server error from the wrapped error.
package apperr

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Err is a interface that represents an error.
type Err interface {
	Error() string
	Code() codes.Code
	DetailCode() string
	Message() string
	Log() string
	Type() Type

	HTTPStatus() int
	GRPCError(domain string) error
}

// Extract is a function to extract apperr.Err from an error.
func Extract(err error) (Err, bool) {
	var e Err
	if ok := errors.As(err, &e); ok {
		return e, true
	}
	return nil, false
}

// ExtractFromGRPCError is a function to extract apperr.Err from a gRPC error.
func ExtractFromGRPCError(err error) (Err, bool) {

	sts, ok := status.FromError(err)
	if !ok {
		return nil, false
	}

	var detailCode string
	d := sts.Details()
	if len(d) > 0 {
		if e, ok := d[0].(*errdetails.ErrorInfo); ok {
			detailCode = e.GetReason()
		}
	}

	switch sts.Code() {
	case codes.PermissionDenied, codes.Unauthenticated,
		codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.FailedPrecondition,
		codes.Canceled, codes.ResourceExhausted, codes.Aborted, codes.OutOfRange:
		return NewClientError(sts.Code(), detailCode, sts.Message()), true
	case codes.Internal, codes.Unavailable, codes.Unimplemented,
		codes.DeadlineExceeded, codes.DataLoss, codes.Unknown:
		return NewServerError(sts.Code(), detailCode, sts.Message(), "grpc error is a server error"), true
	default: // codes.OK
		return nil, false
	}
}
