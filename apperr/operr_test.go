package apperr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/takuoki/golib/apperr"
)

func TestError(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		err := apperr.Err{
			Status:  1,
			Code:    "code",
			Message: "message",
			Log:     "log",
		}
		err2 := fmt.Errorf("wrapped: %w", &err)
		err3 := fmt.Errorf("wrapped: %w", err2)
		result, ok := apperr.Extract(err3)
		if result == nil {
			t.Error("result value must not be nil")
		} else if *result != err {
			t.Errorf("result value does not match expected value (want=%+v, actual=%+v)", err, *result)
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
