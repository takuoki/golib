// Package echoctx provides a custom context for echo.
// Primarily in middleware, use this custom context to set values for context.
// See https://echo.labstack.com/guide/context/ for more details.
package echoctx

import (
	"context"

	"github.com/labstack/echo/v4"
)

// Context is a custom context for echo.
type Context interface {
	echo.Context
	SetContext(ctx context.Context) Context
	GetContext() context.Context
}

type customContext struct {
	echo.Context
	ctx context.Context
}

// New returns new custom context.
// If the argument context is already a custom context, it is returned as is.
func New(c echo.Context) Context {
	if cc, ok := c.(Context); ok {
		return cc
	}
	return &customContext{Context: c, ctx: c.Request().Context()}
}

// GetContext returns pure context.
// If you want to set a value to a context,
// set the value to the context returned by this method.
func (c *customContext) GetContext() context.Context {
	return c.ctx
}

// SetContext sets the context to a custom context.
// If you change the context, be sure to use this method to set the context.
func (c *customContext) SetContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}
