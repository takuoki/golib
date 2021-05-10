package apperr

import "errors"

// Err is a struct that represents an error.
type Err struct {
	Status  int
	Code    string
	Message string
	Log     string
}

// Error is a method to satisfy the error interface.
func (e *Err) Error() string {
	return e.Message
}

// Extract is a function to extract operr.Err from an error.
func Extract(err error) (*Err, bool) {
	e := &Err{}
	if ok := errors.As(err, &e); !ok {
		return nil, false
	}
	return e, true
}
