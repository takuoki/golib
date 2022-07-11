package apperr

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverError struct {
	code       codes.Code
	detailCode string
	message    string
	log        string
}

// NewServerError creates new server error.
// Set the gRPC code to `code`.
// Set the error detail code (ex. "S0001") that the client can handle to `detailCode`.
func NewServerError(code codes.Code, detailCode, message, log string) Err {
	return &serverError{
		code:       code,
		detailCode: detailCode,
		message:    message,
		log:        log,
	}
}

// Error is a method to satisfy the error interface.
func (e *serverError) Error() string {
	return e.message
}

// Code returns code value.
func (e *serverError) Code() codes.Code {
	return e.code
}

// DetailCode returns detail code string.
func (e *serverError) DetailCode() string {
	return e.detailCode
}

// Message returns message string.
func (e *serverError) Message() string {
	return e.message
}

// Log returns log string.
func (e *serverError) Log() string {
	return e.log
}

// Type returns error type.
func (e *serverError) Type() Type {
	return ServerError
}

// GRPCError returns gRPC error.
func (e *serverError) GRPCError(domain string) error {
	st := status.New(e.Code(), e.Message())
	st, _ = st.WithDetails(&errdetails.ErrorInfo{
		Reason: e.DetailCode(),
		Domain: domain,
	})

	return st.Err()
}
