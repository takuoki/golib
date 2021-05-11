package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/takuoki/golib/apperr"
)

func TestErr(t *testing.T) {
	t.Run("client", func(t *testing.T) {
		err := apperr.ClientError(1, "code", "message")
		if err.Error() != "message" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "message", err.Error())
		}
		if err.Status() != 1 {
			t.Errorf("result value does not match expected value (want=%v, actual=%v)", 1, err.Status())
		}
		if err.Code() != "code" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "code", err.Code())
		}
		if err.Message() != "message" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "message", err.Message())
		}
		if err.Log() != "" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "", err.Log())
		}
	})
	t.Run("server", func(t *testing.T) {
		err := apperr.ServerError(1, "code", "message", "log")
		if err.Error() != "message" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "message", err.Error())
		}
		if err.Status() != 1 {
			t.Errorf("result value does not match expected value (want=%v, actual=%v)", 1, err.Status())
		}
		if err.Code() != "code" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "code", err.Code())
		}
		if err.Message() != "message" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "message", err.Message())
		}
		if err.Log() != "log" {
			t.Errorf("result value does not match expected value (want=%q, actual=%q)", "log", err.Log())
		}
	})
}

func TestExtract(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		err := apperr.ServerError(1, "code", "message", "log")
		err2 := fmt.Errorf("wrapped: %w", err)
		err3 := fmt.Errorf("wrapped: %w", err2)
		result, ok := apperr.Extract(err3)
		if result == nil {
			t.Error("result value must not be nil")
		} else if result != err {
			t.Errorf("result value does not match expected value (want=%+v, actual=%+v)", err, result)
		}
		if !ok {
			t.Error("'ok' value must be true")
		}
	})
	t.Run("not-exist", func(t *testing.T) {
		err := errors.New("error")
		result, ok := apperr.Extract(err)
		if result != nil {
			t.Error("result value must be nil")
		}
		if ok {
			t.Error("'ok' value must be false")
		}
	})
	t.Run("nil", func(t *testing.T) {
		result, ok := apperr.Extract(nil)
		if result != nil {
			t.Error("result value must be nil")
		}
		if ok {
			t.Error("'ok' value must be false")
		}
	})
}
