package echoctx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/takuoki/golib/appctx/echoctx"
)

// nolint:staticcheck
func TestContext(t *testing.T) {

	middlewareFunc := func(key, value string) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				ec := echoctx.New(c)
				ctx := ec.GetContext()
				ctx = context.WithValue(ctx, key, value)
				return next(ec.SetContext(ctx))
			}
		}
	}

	testcases := map[string]struct {
		middlewares   []echo.MiddlewareFunc
		wantKeyValues map[string]string
	}{
		"no middleware": {
			middlewares:   nil,
			wantKeyValues: nil,
		},
		"1 middleware": {
			middlewares: []echo.MiddlewareFunc{
				middlewareFunc("key", "value"),
			},
			wantKeyValues: map[string]string{
				"key": "value",
			},
		},
		"2 middlewares": {
			middlewares: []echo.MiddlewareFunc{
				middlewareFunc("key1", "value1"),
				middlewareFunc("key2", "value2"),
			},
			wantKeyValues: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			h := func(c echo.Context) error {
				ctx := echoctx.New(c).GetContext()
				for k, v := range tc.wantKeyValues {
					assert.Equal(t, v, ctx.Value(k))
				}
				return c.NoContent(http.StatusOK)
			}

			for i := len(tc.middlewares) - 1; 0 <= i; i-- {
				h = tc.middlewares[i](h)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			c := e.NewContext(req, rec)

			err := h(c)
			assert.NoError(t, err)
		})
	}
}
