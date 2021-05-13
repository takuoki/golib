package applog_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/applog"
)

func TestBasicLoggerPrint(t *testing.T) {
	defaultTimeFormat := "15:04:05"
	testcase := map[string]struct {
		settingLevel applog.Level
		timeFormat   string
		imageTag     string
		printLevel   applog.Level
		message      string
		labels       map[string]string
		requestID    string
		want         string
	}{
		"default": {
			printLevel: applog.InfoLevel,
			message:    "message",
			want:       `{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"message"}` + "\n",
		},
		"time-format": {
			timeFormat: "2006/01/02 15:04:05",
			printLevel: applog.InfoLevel,
			message:    "message",
			want:       `{"time":"\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}","level":"INFO","message":"message"}` + "\n",
		},
		"image-tag": {
			imageTag:   "image-tag",
			printLevel: applog.InfoLevel,
			message:    "message",
			want:       `{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"message","image_tag":"image-tag"}` + "\n",
		},
		"request-id": {
			requestID:  "request-id",
			printLevel: applog.InfoLevel,
			message:    "message",
			want:       `{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"message","request_id":"request-id"}` + "\n",
		},
		"labels": {
			labels:     map[string]string{"foo": "abc", "bar": "xyz"},
			printLevel: applog.InfoLevel,
			message:    "message",
			want:       `{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"message","labels":{"bar":"xyz","foo":"abc"}}` + "\n",
		},
		"break-line": {
			printLevel: applog.InfoLevel,
			message:    "message\nmessage",
			want:       `{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"message\\nmessage"}` + "\n",
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
			if c.timeFormat != "" {
				opts = append(opts, applog.TimeFormatOption(c.timeFormat))
			} else {
				opts = append(opts, applog.TimeFormatOption(defaultTimeFormat))
			}
			if c.imageTag != "" {
				opts = append(opts, applog.ImageTagOption(c.imageTag))
			}

			buf := &bytes.Buffer{}
			logger := applog.NewBasicLogger(buf, opts...)

			ctx := context.Background()
			if c.requestID != "" {
				ctx = appctx.WithRequestID(ctx, c.requestID)
			}
			logger.Print(ctx, c.printLevel, c.message, c.labels)

			assert.Regexp(t, "^"+c.want+"$", buf.String())
		})
	}
}

func TestBasicLoggerLevel(t *testing.T) {
	newLogger := func(w io.Writer) applog.Logger {
		return applog.NewBasicLogger(
			w,
			applog.TimeFormatOption("15:04:05"),
			applog.LevelOption(applog.UnknownLevel))
	}
	t.Run("critical", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Critical(context.Background(), "message")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"CRITICAL","message":"message"}`+"\n$", buf.String())
	})
	t.Run("error", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Error(context.Background(), "message")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"ERROR","message":"message"}`+"\n$", buf.String())
	})
	t.Run("warn", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Warn(context.Background(), "message")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"WARN","message":"message"}`+"\n$", buf.String())
	})
	t.Run("info", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Info(context.Background(), "message")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"message"}`+"\n$", buf.String())
	})
	t.Run("debug", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Debug(context.Background(), "message")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"DEBUG","message":"message"}`+"\n$", buf.String())
	})
	t.Run("trace", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Trace(context.Background(), "message")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"TRACE","message":"message"}`+"\n$", buf.String())
	})
	t.Run("criticalf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Criticalf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"CRITICAL","message":"value: abc"}`+"\n$", buf.String())
	})
	t.Run("errorf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Errorf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"ERROR","message":"value: abc"}`+"\n$", buf.String())
	})
	t.Run("warnf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Warnf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"WARN","message":"value: abc"}`+"\n$", buf.String())
	})
	t.Run("infof", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Infof(context.Background(), "value: %s", "abc")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"INFO","message":"value: abc"}`+"\n$", buf.String())
	})
	t.Run("debugf", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Debugf(context.Background(), "value: %s", "abc")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"DEBUG","message":"value: abc"}`+"\n$", buf.String())
	})
	t.Run("tracef", func(t *testing.T) {
		buf := &bytes.Buffer{}
		newLogger(buf).Tracef(context.Background(), "value: %s", "abc")
		assert.Regexp(t, `^{"time":"\d{2}:\d{2}:\d{2}","level":"TRACE","message":"value: abc"}`+"\n$", buf.String())
	})
}
