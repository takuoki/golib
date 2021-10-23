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
		assert.Equal(t, testRequestID, result, "RequestID is not equal")
	})
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		result := appctx.RequestID(ctx)
		assert.Empty(t, result, "RequestID is not empty")
	})
}

func TestUserID(t *testing.T) {
	t.Run("succsess", func(t *testing.T) {
		testUserID := "test-user-id"
		ctx := context.Background()
		ctx = appctx.WithUserID(ctx, testUserID)
		result := appctx.UserID(ctx)
		assert.Equal(t, testUserID, result, "UserID is not equal")
	})
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		result := appctx.UserID(ctx)
		assert.Empty(t, result, "UserID is not empty")
	})
}

func TestAuthorization(t *testing.T) {
	t.Run("succsess", func(t *testing.T) {
		testAuthorization := "test-authorization"
		ctx := context.Background()
		ctx = appctx.WithAuthorization(ctx, testAuthorization)
		result := appctx.Authorization(ctx)
		assert.Equal(t, testAuthorization, result, "Authorization is not equal")
	})
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		result := appctx.Authorization(ctx)
		assert.Empty(t, result, "Authorization is not empty")
	})
}
