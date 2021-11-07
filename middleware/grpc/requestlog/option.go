package grpc_requestlog

import "github.com/google/uuid"

type options struct {
	requestIDKey  string
	requestIDFunc func() (string, error)
}

var defaultOptions = options{
	requestIDKey:  "Request-ID",
	requestIDFunc: defaultRequestIDFunc,
}

func defaultRequestIDFunc() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Option is an option when creating middleware.
type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// RequestIDKey is a key option for the request ID specified in the metadata.
// The default is "Request-ID".
func RequestIDKey(key string) Option {
	return newFuncOption(func(o *options) {
		o.requestIDKey = key
	})
}

// RequestIDFunc is a function option to automatically generate a request ID
// when it is not specified. The default is a function that generates a UUID.
func RequestIDFunc(fn func() (string, error)) Option {
	return newFuncOption(func(o *options) {
		o.requestIDFunc = fn
	})
}
