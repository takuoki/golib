// Package apperr defines an error that can distinguish
// between a client error and a server error.
// It also defines a function that extracts the defined
// client or server error from the wrapped error.
package apperr

import "errors"

// Err is a interface that represents an error.
type Err interface {
	Error() string
	Status() int
	Code() string
	Message() string
	Log() string
	IsClientError() bool
	IsServerError() bool
}

type er struct {
	status   int
	code     string
	message  string
	log      string
	isClient bool
}

// ClientError creates new client error.
func ClientError(status int, code, message string) Err {
	return &er{
		status:   status,
		code:     code,
		message:  message,
		isClient: true,
	}
}

// ServerError creates new server error.
func ServerError(status int, code, message, log string) Err {
	return &er{
		status:   status,
		code:     code,
		message:  message,
		log:      log,
		isClient: false,
	}
}

// Error is a method to satisfy the error interface.
func (e *er) Error() string {
	return e.message
}

// Status returns status value.
func (e *er) Status() int {
	return e.status
}

// Code returns code string.
func (e *er) Code() string {
	return e.code
}

// Message returns message string.
func (e *er) Message() string {
	return e.message
}

// Log returns log string.
func (e *er) Log() string {
	return e.log
}

// IsClientError returns whether it is a client error.
func (e *er) IsClientError() bool {
	return e.isClient
}

// IsServerError returns whether it is a server error.
func (e *er) IsServerError() bool {
	return !e.isClient
}

// Extract is a function to extract apperr.Err from an error.
func Extract(err error) (Err, bool) {
	e := &er{}
	if ok := errors.As(err, &e); !ok {
		return nil, false
	}
	return e, true
}
