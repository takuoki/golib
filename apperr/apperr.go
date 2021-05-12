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

type cerr struct {
	status  int
	code    string
	message string
}

// ClientError creates new client error.
// For RESTful API, set the HTTP status code to `status`.
// Set the error code (ex. "E0001") that the client can handle to `code`.
func ClientError(status int, code, message string) Err {
	return &cerr{
		status:  status,
		code:    code,
		message: message,
	}
}

// Error is a method to satisfy the error interface.
func (e *cerr) Error() string {
	return e.message
}

// Status returns status value.
func (e *cerr) Status() int {
	return e.status
}

// Code returns code string.
func (e *cerr) Code() string {
	return e.code
}

// Message returns message string.
func (e *cerr) Message() string {
	return e.message
}

// Log returns log string.
func (e *cerr) Log() string {
	return ""
}

// IsClientError returns whether it is a client error.
func (e *cerr) IsClientError() bool {
	return true
}

// IsServerError returns whether it is a server error.
func (e *cerr) IsServerError() bool {
	return false
}

type serr struct {
	status  int
	code    string
	message string
	log     string
}

// ServerError creates new server error.
// For RESTful API, set the HTTP status code to `status`.
// Set the error code (ex. "S0001") that the client can handle to `code`.
func ServerError(status int, code, message, log string) Err {
	return &serr{
		status:  status,
		code:    code,
		message: message,
		log:     log,
	}
}

// Error is a method to satisfy the error interface.
func (e *serr) Error() string {
	return e.message
}

// Status returns status value.
func (e *serr) Status() int {
	return e.status
}

// Code returns code string.
func (e *serr) Code() string {
	return e.code
}

// Message returns message string.
func (e *serr) Message() string {
	return e.message
}

// Log returns log string.
func (e *serr) Log() string {
	return e.log
}

// IsClientError returns whether it is a client error.
func (e *serr) IsClientError() bool {
	return false
}

// IsServerError returns whether it is a server error.
func (e *serr) IsServerError() bool {
	return true
}

// Extract is a function to extract apperr.Err from an error.
func Extract(err error) (Err, bool) {
	var e Err
	if ok := errors.As(err, &e); ok {
		return e, true
	}
	return nil, false
}
