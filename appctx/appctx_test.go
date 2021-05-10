package appctx_test

import (
	"context"
	"testing"

	"github.com/takuoki/golib/appctx"
)

func TestRequestID(t *testing.T) {
	t.Run("succsess", func(t *testing.T) {
		testRequestID := "test-request-id"
		ctx := context.Background()
		ctx = appctx.SetRequestID(ctx, testRequestID)
		result := appctx.GetRequestID(ctx)
		if result != testRequestID {
			t.Errorf("value does not match the expected value (want=%q, actual=%q)", testRequestID, result)
		}
	})
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		result := appctx.GetRequestID(ctx)
		if result != "" {
			t.Errorf("value must be empty string (actual=%q)", result)
		}
	})
}
