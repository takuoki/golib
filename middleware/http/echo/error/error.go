// Package echo_error is a generic server-side echo middleware
// that converts standard error to HTTP error.
package echo_error

import (
	"context"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"

	"github.com/takuoki/golib/appctx/echoctx"
	"github.com/takuoki/golib/apperr"
	"github.com/takuoki/golib/applog"
)

// Middleware returns a echo middleware that converts standard error to HTTP error.
func Middleware(internalServerErrorCode string, logger applog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				e, ok := apperr.Extract(err)
				if !ok {
					if herr, ok := err.(*echo.HTTPError); ok {
						e = apperr.NewClientError(
							codeFromHTTPStatus(echoctx.New(c).GetContext(), herr.Code, logger),
							"-",
							fmt.Sprintf("%v", herr.Message),
						)
					} else {
						e = newInternalServerError(internalServerErrorCode, err)
					}
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

func codeFromHTTPStatus(ctx context.Context, status int, logger applog.Logger) codes.Code {
	switch status {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return codes.OK
	case http.StatusRequestTimeout:
		return codes.Canceled
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	}

	logger.Warnf(ctx, "unknown HTTP status: %d", status)
	return codes.Internal
}
