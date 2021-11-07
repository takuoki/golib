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

func RequestIDKey(key string) Option {
	return newFuncOption(func(o *options) {
		o.requestIDKey = key
	})
}

func RequestIDFunc(fn func() (string, error)) Option {
	return newFuncOption(func(o *options) {
		o.requestIDFunc = fn
	})
}
