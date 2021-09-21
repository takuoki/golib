package applog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/takuoki/golib/appctx"
)

type basicLogger struct {
	mu         sync.Mutex
	out        io.Writer
	level      Level
	timeFormat string
	imageTag   string
}

// NewBasicLogger creates a basic logger that outputs in JSON format.
// It handles the necessity of output according to the log level,
// and outputs context information to the log in common.
func NewBasicLogger(w io.Writer, opts ...Option) Logger {
	logger := &basicLogger{
		out:        w,
		timeFormat: time.RFC3339,
	}
	for _, opt := range opts {
		// basicLogger option never returns an error
		_ = opt(logger)
	}
	return logger
}

func (l *basicLogger) SetLevel(lv Level) error {
	l.level = lv
	return nil
}

func (l *basicLogger) SetTimeFormat(format string) error {
	l.timeFormat = format
	return nil
}

func (l *basicLogger) SetImageTag(tag string) error {
	l.imageTag = tag
	return nil
}

func (l *basicLogger) Critical(ctx context.Context, msg string) {
	l.Print(ctx, CriticalLevel, msg, nil)
}

func (l *basicLogger) Error(ctx context.Context, msg string) {
	l.Print(ctx, ErrorLevel, msg, nil)
}

func (l *basicLogger) Warn(ctx context.Context, msg string) {
	l.Print(ctx, WarnLevel, msg, nil)
}

func (l *basicLogger) Info(ctx context.Context, msg string) {
	l.Print(ctx, InfoLevel, msg, nil)
}

func (l *basicLogger) Debug(ctx context.Context, msg string) {
	l.Print(ctx, DebugLevel, msg, nil)
}

func (l *basicLogger) Trace(ctx context.Context, msg string) {
	l.Print(ctx, TraceLevel, msg, nil)
}

func (l *basicLogger) Criticalf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, CriticalLevel, format, a...)
}

func (l *basicLogger) Errorf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, ErrorLevel, format, a...)
}

func (l *basicLogger) Warnf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, WarnLevel, format, a...)
}

func (l *basicLogger) Infof(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, InfoLevel, format, a...)
}

func (l *basicLogger) Debugf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, DebugLevel, format, a...)
}

func (l *basicLogger) Tracef(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, TraceLevel, format, a...)
}

func (l *basicLogger) printf(ctx context.Context, lv Level, format string, a ...interface{}) {
	l.Print(ctx, lv, fmt.Sprintf(format, a...), nil)
}

func (l *basicLogger) Print(ctx context.Context, lv Level, msg string, labels map[string]string) {
	if !shouldPrint(l.level, lv) {
		return
	}
	log := basicLog{
		Time:      time.Now().Format(l.timeFormat),
		Level:     lv.String(),
		Message:   msg,
		ImageTag:  l.imageTag,
		RequestID: appctx.RequestID(ctx),
		Labels:    labels,
	}
	jsonLog, _ := json.Marshal(log)

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, string(jsonLog))
}

type basicLog struct {
	Time      string            `json:"time"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	ImageTag  string            `json:"image_tag,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}
