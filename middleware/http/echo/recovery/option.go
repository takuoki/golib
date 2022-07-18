package echo_recovery

import (
	"context"
)

type options struct {
	recoveryFunc func(ctx context.Context, p interface{}) (err error)
}

var defaultOptions = options{
	recoveryFunc: nil,
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

// RecoveryFunc is a function option for recovering from a panic.
// Be sure to specify either RecoveryFunc or RecoveryContextFunc.
func RecoveryFunc(fn func(p interface{}) (err error)) Option {
	return RecoveryContextFunc(func(ctx context.Context, p interface{}) (err error) {
		return fn(p)
	})
}

// RecoveryContextFunc is a function option for recovering from a panic.
// Be sure to specify either RecoveryFunc or RecoveryContextFunc.
func RecoveryContextFunc(fn func(ctx context.Context, p interface{}) (err error)) Option {
	return newFuncOption(func(o *options) {
		o.recoveryFunc = fn
	})
}
