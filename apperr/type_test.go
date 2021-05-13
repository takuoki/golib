package apperr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/apperr"
)

func TestType(t *testing.T) {
	testcase := map[string]struct {
		in   apperr.Type
		want string
	}{
		"client":  {in: apperr.ClientError, want: "client-error"},
		"server":  {in: apperr.ServerError, want: "server-error"},
		"unknown": {in: 0, want: "unknown"},
	}

	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			r := c.in.String()
			assert.Equal(t, c.want, r)
		})
	}
}
