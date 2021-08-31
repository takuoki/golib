package appnotice_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/appnotice"
)

func TestNopNotifier(t *testing.T) {
	n := appnotice.NewNopNotifier()
	assert.Nil(t, n.Error(errors.New("error")), "response of Error must be nil")
	assert.Nil(t, n.Critical(errors.New("critical")), "response of Critical must be nil")
}
