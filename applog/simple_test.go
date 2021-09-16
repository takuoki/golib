package applog_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/applog"
)

func TestSimpleLoggerPrint(t *testing.T) {
	testcase := map[string]struct {
		settingLevel applog.Level
		printLevel   applog.Level
		message      string
		labels       map[string]string
		want         string
	}{
		"default": {
			printLevel: applog.InfoLevel,
			message:    "message",
			want:       "message\n",
		},
		"labels": {
			printLevel: applog.InfoLevel,
			message:    "message",
			labels:     map[string]string{"foo": "abc", "bar": "xyz"},
			want:       `message \(bar: xyz, foo: abc\)` + "\n",
		},
		"break-line": {
			printLevel: applog.InfoLevel,
			message:    "message\nmessage",
			want:       "message\nmessage\n",
		},
		"no-output": {
			settingLevel: applog.ErrorLevel,
			printLevel:   applog.InfoLevel,
			message:      "message",
			want:         "",
		},
	}

	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			opts := []applog.Option{}
			if c.settingLevel != 0 {
				opts = append(opts, applog.LevelOption(c.settingLevel))
			}

			buf := &bytes.Buffer{}
			logger, err := applog.NewSimpleLogger(buf, opts...)
			if err != nil {
				t.Fatalf("error occurred in NewSimpleLogger: %v", err)
			}

			ctx := context.Background()
			logger.Print(ctx, c.printLevel, c.message, c.labels)

			assert.Regexp(t, "^"+c.want+"$", buf.String())
		})
	}
}

func TestSimpleLoggerLevel(t *testing.T) {
	newLogger := func(w io.Writer) applog.Logger {
		l, err := applog.NewSimpleLogger(w, applog.LevelOption(applog.UnknownLevel))
		if err != nil {
			t.Fatalf("error occurred in NewSimpleLogger: %v", err)
		}
		return l
	}
	t.Run("critical", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Critical(context.Background(), "message")
		assert.Regexp(t, "^message\n$", buf.String())
	})
	t.Run("error", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Error(context.Background(), "message")
		assert.Regexp(t, "^message\n$", buf.String())
	})
	t.Run("warn", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Warn(context.Background(), "message")
		assert.Regexp(t, "^message\n$", buf.String())
	})
	t.Run("info", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Info(context.Background(), "message")
		assert.Regexp(t, "^message\n$", buf.String())
	})
	t.Run("debug", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Debug(context.Background(), "message")
		assert.Regexp(t, "^message\n$", buf.String())
	})
	t.Run("trace", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Trace(context.Background(), "message")
		assert.Regexp(t, "^message\n$", buf.String())
	})
	t.Run("criticalf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Criticalf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, "^value: abc\n$", buf.String())
	})
	t.Run("errorf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Errorf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, "^value: abc\n$", buf.String())
	})
	t.Run("warnf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Warnf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, "^value: abc\n$", buf.String())
	})
	t.Run("infof", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Infof(context.Background(), "value: %s", "abc")
		assert.Regexp(t, "^value: abc\n$", buf.String())
	})
	t.Run("debugf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Debugf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, "^value: abc\n$", buf.String())
	})
	t.Run("tracef", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Tracef(context.Background(), "value: %s", "abc")
		assert.Regexp(t, "^value: abc\n$", buf.String())
	})
}

func TestSimpleLoggerOptionError(t *testing.T) {
	t.Run("timeFormat", func(t *testing.T) {
		buf := &bytes.Buffer{}
		_, err := applog.NewSimpleLogger(buf, applog.TimeFormatOption("dummy"))
		if assert.NotNil(t, err) {
			assert.Equal(t, "TimeFormatOption is not available for simpleLogger", err.Error())
		}
	})
	t.Run("imageTag", func(t *testing.T) {
		buf := &bytes.Buffer{}
		_, err := applog.NewSimpleLogger(buf, applog.ImageTagOption("dummy"))
		if assert.NotNil(t, err) {
			assert.Equal(t, "ImageTagOption is not available for simpleLogger", err.Error())
		}
	})
}
