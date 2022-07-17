// Package echo_requestlog is a generic server-side echo middleware
// that outputs request log.
// If the request ID is not specified, it will be automatically generated.
package echo_requestlog

import (
	"strings"

	echo "github.com/labstack/echo/v4"

	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/appctx/echoctx"
	"github.com/takuoki/golib/applog"
)

// Middleware returns a echo middleware that outputs request logs.
func Middleware(logger applog.Logger, opt ...Option) echo.MiddlewareFunc {

	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ec := echoctx.New(c)
			ctx := ec.GetContext()

			// Set request ID to context.
			reqID := c.Request().Header.Get(opts.requestIDKey)
			if reqID == "" {
				id, err := opts.requestIDFunc()
				if err != nil {
					logger.Warnf(ctx, "failed to create new request ID: %v", err)
				}
				reqID = id
			}
			ctx = appctx.WithRequestID(ctx, reqID)

			// Create label.
			label := map[string]string{
				"host":   c.Request().Host,
				"method": c.Request().Method,
				"uri":    c.Request().RequestURI,
			}

			ip := c.Request().RemoteAddr
			if idx := strings.LastIndex(ip, ":"); idx > 0 {
				ip = ip[0:idx]
			}
			if ip != "" {
				label["ip_address"] = ip
			}
			if ua := c.Request().UserAgent(); ua != "" {
				label["user_agent"] = ua
			}

			logger.Print(ctx, applog.InfoLevel, "request log", label)

			return next(ec.SetContext(ctx))
		}
	}
}
