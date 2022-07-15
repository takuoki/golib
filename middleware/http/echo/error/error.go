// Package echo_error is a generic server-side echo middleware
// that converts standard error to HTTP error.
package echo_error

import (
	"google.golang.org/grpc/codes"

	echo "github.com/labstack/echo/v4"
	"github.com/takuoki/golib/apperr"
	"github.com/takuoki/golib/applog"
)

func Middleware(internalServerErrorCode string, logger applog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				e, ok := apperr.Extract(err)
				if !ok {
					e = newInternalServerError(internalServerErrorCode, err)
				}
				if e.Log() != "" {
					logger.Error(c.Request().Context(), e.Log())
				}

				return c.JSON(e.HTTPStatus(), struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    e.DetailCode(),
					Message: e.Message(),
				})
			}

			return nil
		}
	}
}

func newInternalServerError(code string, err error) apperr.Err {
	return apperr.NewServerError(
		codes.Internal,
		code,
		"internal server error",
		err.Error(),
	)
}
