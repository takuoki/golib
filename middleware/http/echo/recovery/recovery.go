// Package echo_recovery is a generic server-side echo middleware
// that recovers panic.
package echo_recovery

import (
	"context"
	"fmt"

	echo "github.com/labstack/echo/v4"

	"github.com/takuoki/golib/appctx/echoctx"
)

// Middleware returns a echo middleware that recovers panic.
func Middleware(opt ...Option) echo.MiddlewareFunc {

	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = recoverFrom(echoctx.New(c).GetContext(), r, opts.recoveryFunc)
				}
			}()

			return next(c)
		}
	}
}

func recoverFrom(ctx context.Context, p interface{}, r func(ctx context.Context, p interface{}) (err error)) error {
	if r == nil {
		return fmt.Errorf("recovery function is nil: %v", p)
	}
	return r(ctx, p)
}
