package recovery_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/recovery"
)

func basicUsage(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recovery.Recovery(r)
		}
	}()

	er := f()
	return er
}

func TestRecovery(t *testing.T) {
	testcase := map[string]struct {
		f       func() error
		wantErr string
	}{
		"no-error":  {f: func() error { return nil }},
		"error":     {f: func() error { return errors.New("error message") }, wantErr: "error message"},
		"panic":     {f: func() error { panic("panic message") }, wantErr: "panic recovered: panic message"},
		"nil-panic": {f: func() error { panic(nil) }, wantErr: "panic recovered: panic called with nil argument"},
	}

	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			err := basicUsage(c.f)
			if c.wantErr == "" {
				assert.Nil(t, err)
			} else {
				assert.Regexp(t, "^"+c.wantErr, err.Error())
			}
		})
	}
}
