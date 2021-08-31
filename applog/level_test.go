package applog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/applog"
)

func TestParseLevel(t *testing.T) {
	testcase := map[string]struct {
		in   string
		want applog.Level
		err  bool
	}{
		"critical": {in: "critical", want: applog.CriticalLevel},
		"error":    {in: "error", want: applog.ErrorLevel},
		"warn":     {in: "warn", want: applog.WarnLevel},
		"info":     {in: "info", want: applog.InfoLevel},
		"debug":    {in: "debug", want: applog.DebugLevel},
		"trace":    {in: "trace", want: applog.TraceLevel},
		"empty":    {in: "", want: applog.UnknownLevel, err: true},
	}

	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			lv, err := applog.ParseLevel(c.in)
			assert.Equal(t, c.want, lv)
			if c.err {
				assert.NotNil(t, err, "err is nil")
			} else {
				assert.Nil(t, err, "err is not nil")
			}
		})
	}
}

func TestParseLevelWithDefault(t *testing.T) {
	testcase := map[string]struct {
		in   string
		def  applog.Level
		want applog.Level
	}{
		"critical": {in: "critical", def: applog.DebugLevel, want: applog.CriticalLevel},
		"error":    {in: "error", def: applog.DebugLevel, want: applog.ErrorLevel},
		"empty":    {in: "", def: applog.DebugLevel, want: applog.DebugLevel},
	}

	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.want, applog.ParseLevelWithDefault(c.in, c.def))
		})
	}
}

func TestLevel(t *testing.T) {
	testcase := map[string]struct {
		in   applog.Level
		want string
	}{
		"critical": {in: applog.CriticalLevel, want: "CRITICAL"},
		"error":    {in: applog.ErrorLevel, want: "ERROR"},
		"warn":     {in: applog.WarnLevel, want: "WARN"},
		"info":     {in: applog.InfoLevel, want: "INFO"},
		"debug":    {in: applog.DebugLevel, want: "DEBUG"},
		"trace":    {in: applog.TraceLevel, want: "TRACE"},
		"unknown":  {in: applog.UnknownLevel, want: "UNKNOWN"},
	}

	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.want, c.in.String())
		})
	}
}
