package echo_requestlog_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/appctx/echoctx"
	"github.com/takuoki/golib/applog"
	echo_requestlog "github.com/takuoki/golib/middleware/http/echo/requestlog"
)

func TestMiddleware(t *testing.T) {

	const requestIDKey = "Request-ID"
	const originKey = "Origin"
	const userAgentKey = "User-Agent"

	testcases := map[string]struct {
		opts      []echo_requestlog.Option
		method    string
		uri       string
		reqID     string
		origin    string
		userAgent string
		wantLog   string
		wantReqID string
	}{
		"empty reqID and userAgent": {
			opts: []echo_requestlog.Option{
				echo_requestlog.RequestIDFunc(func() (string, error) { return "req-id", nil }),
			},
			method:    http.MethodGet,
			uri:       "/test?a=1&b=2#xyz",
			reqID:     "",
			origin:    "",
			userAgent: "",
			wantLog:   `^request log \(host: .*, ip_address: [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+, method: GET, uri: /test\?a=1&b=2#xyz\)` + "\n$",
			wantReqID: "req-id",
		},
		"exist reqID, origin and userAgent": {
			opts: []echo_requestlog.Option{
				echo_requestlog.RequestIDFunc(func() (string, error) { return "new-req-id", nil }),
			},
			method:    http.MethodPost,
			uri:       "/test",
			origin:    "http://localhost:8080",
			reqID:     "req-id",
			userAgent: "user-agent",
			wantLog:   `^request log \(host: .*, ip_address: [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+, method: POST, origin: http://localhost:8080, uri: /test, user_agent: user-agent\)` + "\n$",
			wantReqID: "req-id",
		},
		"create reqID error": {
			opts: []echo_requestlog.Option{
				echo_requestlog.RequestIDFunc(func() (string, error) { return "", errors.New("error") }),
			},
			method: http.MethodGet,
			uri:    "/test",
			wantLog: "^failed to create new request ID: error\n" +
				`request log \(host: .*, ip_address: [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+, method: GET, uri: /test\)` + "\n$",
			wantReqID: "",
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

			m := echo_requestlog.Middleware(logger, tc.opts...)
			h := m(func(c echo.Context) error {
				assert.Equal(t, tc.wantReqID, appctx.RequestID(echoctx.New(c).GetContext()))
				return c.NoContent(http.StatusOK)
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.uri, nil)
			if tc.reqID != "" {
				req.Header.Add(requestIDKey, tc.reqID)
			}
			if tc.origin != "" {
				req.Header.Add(originKey, tc.origin)
			}
			if tc.userAgent != "" {
				req.Header.Add(userAgentKey, tc.userAgent)
			}
			c := e.NewContext(req, rec)

			err = h(c)
			assert.NoError(t, err)
			assert.Regexp(t, tc.wantLog, buf.String(), "log doesn't match")
		})
	}
}
