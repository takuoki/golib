package apperr

import "errors"

// Err is a interface that represents an error.
type Err interface {
	Error() string
	Status() int
	Code() string
	Message() string
	Log() string
}

type er struct {
	status  int
	code    string
	message string
	log     string
}

// ClientError creates new client error.
func ClientError(status int, code, message string) Err {
	return &er{
		status:  status,
		code:    code,
		message: message,
	}
}

// ServerError creates new server error.
func ServerError(status int, code, message, log string) Err {
	return &er{
		status:  status,
		code:    code,
		message: message,
		log:     log,
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

// Extract is a function to extract apperr.Err from an error.
func Extract(err error) (Err, bool) {
	e := &er{}
	if ok := errors.As(err, &e); !ok {
		return nil, false
	}
	return e, true
}
