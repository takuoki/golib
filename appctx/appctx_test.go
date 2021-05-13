package appctx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/appctx"
)

func TestRequestID(t *testing.T) {
	t.Run("succsess", func(t *testing.T) {
		testRequestID := "test-request-id"
		ctx := context.Background()
		ctx = appctx.WithRequestID(ctx, testRequestID)
		result := appctx.RequestID(ctx)
		assert.Equal(t, testRequestID, result, "RequestID is not equal.")
	})
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		result := appctx.RequestID(ctx)
		assert.Empty(t, result, "RequestID is not empty.")
	})
}
