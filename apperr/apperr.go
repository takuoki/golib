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
	Type() Type
}

// Extract is a function to extract apperr.Err from an error.
func Extract(err error) (Err, bool) {
	var e Err
	if ok := errors.As(err, &e); ok {
		return e, true
	}
	return nil, false
}
