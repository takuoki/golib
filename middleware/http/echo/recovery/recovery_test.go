package echo_recovery_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/takuoki/golib/appctx/echoctx"
	echo_recovery "github.com/takuoki/golib/middleware/http/echo/recovery"
)

// nolint:staticcheck
func TestMiddleware(t *testing.T) {

	testcases := map[string]struct {
		opts      []echo_recovery.Option
		wantError string
	}{
		"not set": {
			opts:      nil,
			wantError: "recovery function is nil: ",
		},
		"recover": {
			opts: []echo_recovery.Option{
				echo_recovery.RecoveryFunc(func(p interface{}) (err error) {
					return fmt.Errorf("panic recovered: %v", p)
				}),
			},
			wantError: "panic recovered: ",
		},
		"recover context": {
			opts: []echo_recovery.Option{
				echo_recovery.RecoveryContextFunc(func(ctx context.Context, p interface{}) (err error) {
					return fmt.Errorf("panic recovered (%v): %v", ctx.Value("key"), p)
				}),
			},
			wantError: `panic recovered \(value\): `,
		},
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			m := echo_recovery.Middleware(tc.opts...)

			t.Run("success", func(t *testing.T) {
				h := m(func(c echo.Context) error {
					return c.NoContent(http.StatusOK)
				})

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				ec := echoctx.New(e.NewContext(req, rec))
				ec.SetContext(context.WithValue(ec.GetContext(), "key", "value"))

				err := h(ec)
				assert.NoError(t, err)
			})

			t.Run("panic", func(t *testing.T) {
				h := m(func(c echo.Context) error {
					panic("panic")
				})

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				ec := echoctx.New(e.NewContext(req, rec))
				ec.SetContext(context.WithValue(ec.GetContext(), "key", "value"))

				err := h(ec)
				assert.Regexp(t, "^"+tc.wantError, err.Error())
			})
		})
	}
}
