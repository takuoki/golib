package apperr

import "google.golang.org/grpc/codes"

type clientError struct {
	code       codes.Code
	detailCode string
	message    string
}

// NewClientError creates new client error.
// Set the gRPC code to `code`.
// Set the error detail code (ex. "E0001") that the client can handle to `detailCode`.
func NewClientError(code codes.Code, detailCode, message string) Err {
	return &clientError{
		code:       code,
		detailCode: detailCode,
		message:    message,
	}
}

// Error is a method to satisfy the error interface.
func (e *clientError) Error() string {
	return e.message
}

// Code returns code value.
func (e *clientError) Code() codes.Code {
	return e.code
}

// DetailCode returns detail code string.
func (e *clientError) DetailCode() string {
	return e.detailCode
}

// Message returns message string.
func (e *clientError) Message() string {
	return e.message
}

// Log returns log string.
func (e *clientError) Log() string {
	return ""
}

// Type returns error type.
func (e *clientError) Type() Type {
	return ClientError
}
