package notice_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/notice"
)

func TestNopNotifier(t *testing.T) {
	n := notice.NewNopNotifier()
	assert.Nil(t, n.Error(errors.New("error")), "response of Error must be nil")
	assert.Nil(t, n.Critical(errors.New("critical")), "response of Critical must be nil")
}
