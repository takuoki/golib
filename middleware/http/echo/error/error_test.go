package echo_error_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"github.com/takuoki/golib/apperr"
	"github.com/takuoki/golib/applog"
	echo_error "github.com/takuoki/golib/middleware/http/echo/error"
)

func TestAddTrailingSlash(t *testing.T) {

	const internalServerErrorCode = "S0001"

	testcases := map[string]struct {
		err        error
		wantStatus int
		wantResp   string
		wantLog    string
	}{
		"success": {
			err:        nil,
			wantStatus: 200,
			wantResp:   "success",
		},
		"client error": {
			err:        apperr.NewClientError(codes.InvalidArgument, "C0001", "client error"),
			wantStatus: 400,
			wantResp:   `{"code":"C0001","message":"client error"}` + "\n",
		},
		"echo http error": {
			err:        echo.ErrNotFound,
			wantStatus: 404,
			wantResp:   `{"code":"-","message":"Not Found"}` + "\n",
		},
		"server error": {
			err:        errors.New("server error"),
			wantStatus: 500,
			wantResp:   fmt.Sprintf(`{"code":"%s","message":"internal server error"}`+"\n", internalServerErrorCode),
			wantLog:    "server error\n",
		},
	}
	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			e := echo.New()

			buf := &bytes.Buffer{}
			logger, err := applog.NewSimpleLogger(buf)
			if err != nil {
				t.Fatalf("error occurred in NewSimpleLogger: %v", err)
			}

			m := echo_error.Middleware(internalServerErrorCode, logger)
			h := m(func(c echo.Context) error {
				if tc.err != nil {
					return tc.err
				}
				return c.String(http.StatusOK, "success")
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			c := e.NewContext(req, rec)

			err = h(c)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantResp, rec.Body.String())

			if tc.wantLog == "" {
				assert.Empty(t, buf.String(), "log must be empty")
			} else {
				assert.Equal(t, tc.wantLog, buf.String(), "log doesn't match")
			}
		})
	}
}

func TestCodeFromHTTPStatus(t *testing.T) {
	testcases := map[string]struct {
		in      int
		want    codes.Code
		wantLog string
	}{
		"ok":                    {in: http.StatusOK, want: codes.OK},
		"request timeout":       {in: http.StatusRequestTimeout, want: codes.Canceled},
		"bad request":           {in: http.StatusBadRequest, want: codes.InvalidArgument},
		"gateway timeout":       {in: http.StatusGatewayTimeout, want: codes.DeadlineExceeded},
		"not found":             {in: http.StatusNotFound, want: codes.NotFound},
		"conflict":              {in: http.StatusConflict, want: codes.AlreadyExists},
		"forbidden":             {in: http.StatusForbidden, want: codes.PermissionDenied},
		"unauthorized":          {in: http.StatusUnauthorized, want: codes.Unauthenticated},
		"too many requests":     {in: http.StatusTooManyRequests, want: codes.ResourceExhausted},
		"not implemented":       {in: http.StatusNotImplemented, want: codes.Unimplemented},
		"internal server error": {in: http.StatusInternalServerError, want: codes.Internal},
		"service unavailable":   {in: http.StatusServiceUnavailable, want: codes.Unavailable},
		"unknown":               {in: http.StatusTeapot, want: codes.Internal, wantLog: "unknown HTTP status: 418\n"},
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger, err := applog.NewSimpleLogger(buf)
			if err != nil {
				t.Fatalf("error occurred in NewSimpleLogger: %v", err)
			}

			r := echo_error.CodeFromHTTPStatus(context.Background(), tc.in, logger)

			assert.Equal(t, tc.want, r)

			if tc.wantLog == "" {
				assert.Empty(t, buf.String(), "log must be empty")
			} else {
				assert.Equal(t, tc.wantLog, buf.String(), "log doesn't match")
			}
		})
	}
}
