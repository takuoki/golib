package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/apperr"
)

func TestErr(t *testing.T) {
	t.Run("client", func(t *testing.T) {
		err := apperr.NewClientError(1, "code", "message")
		assert.Equal(t, "message", err.Error(), "Error is not equal.")
		assert.Equal(t, 1, err.Status(), "Etatus is not equal.")
		assert.Equal(t, "code", err.Code(), "Code is not equal.")
		assert.Equal(t, "message", err.Message(), "Message is not equal.")
		assert.Equal(t, "", err.Log(), "Log is not equal.")
		assert.Equal(t, apperr.ClientError, err.Type(), "Type is not equal.")
	})
	t.Run("server", func(t *testing.T) {
		err := apperr.NewServerError(1, "code", "message", "log")
		assert.Equal(t, "message", err.Error(), "Error is not equal.")
		assert.Equal(t, 1, err.Status(), "Status is not equal.")
		assert.Equal(t, "code", err.Code(), "Code is not equal.")
		assert.Equal(t, "message", err.Message(), "Message is not equal.")
		assert.Equal(t, "log", err.Log(), "Log is not equal.")
		assert.Equal(t, apperr.ServerError, err.Type(), "Type is not equal.")
	})
}

func TestExtract(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		err := apperr.NewServerError(1, "code", "message", "log")
		err2 := fmt.Errorf("wrapped: %w", err)
		err3 := fmt.Errorf("wrapped: %w", err2)
		result, ok := apperr.Extract(err3)
		if assert.NotNil(t, result, "Error is nil.") {
			assert.Equal(t, err, result, "Error is not equal.")
		}
		assert.True(t, ok, "Ok is not true.")
	})
	t.Run("not-exist", func(t *testing.T) {
		err := errors.New("error")
		result, ok := apperr.Extract(err)
		assert.Nil(t, result, "Error is not nil.")
		assert.False(t, ok, "Ok is not false.")
	})
	t.Run("nil", func(t *testing.T) {
		result, ok := apperr.Extract(nil)
		assert.Nil(t, result, "Error is not nil.")
		assert.False(t, ok, "Ok is not false.")
	})
}
