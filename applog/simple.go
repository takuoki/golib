package applog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

type simpleLogger struct {
	mu    sync.Mutex
	out   io.Writer
	level Level
}

// NewSimpleLogger creates a simple logger that outputs only message.
// It handles the necessity of output according to the log level.
func NewSimpleLogger(w io.Writer, opts ...Option) (Logger, error) {
	logger := &simpleLogger{
		out: w,
	}
	for _, opt := range opts {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}
	return logger, nil
}

func (l *simpleLogger) setLevel(lv Level) error {
	l.level = lv
	return nil
}

func (l *simpleLogger) setTimeFormat(format string) error {
	return errors.New("TimeFormatOption is not available for simpleLogger")
}

func (l *simpleLogger) setImageTag(tag string) error {
	return errors.New("ImageTagOption is not available for simpleLogger")
}

func (l *simpleLogger) Critical(ctx context.Context, msg string) {
	l.Print(ctx, CriticalLevel, msg, nil)
}

func (l *simpleLogger) Error(ctx context.Context, msg string) {
	l.Print(ctx, ErrorLevel, msg, nil)
}

func (l *simpleLogger) Warn(ctx context.Context, msg string) {
	l.Print(ctx, WarnLevel, msg, nil)
}

func (l *simpleLogger) Info(ctx context.Context, msg string) {
	l.Print(ctx, InfoLevel, msg, nil)
}

func (l *simpleLogger) Debug(ctx context.Context, msg string) {
	l.Print(ctx, DebugLevel, msg, nil)
}

func (l *simpleLogger) Trace(ctx context.Context, msg string) {
	l.Print(ctx, TraceLevel, msg, nil)
}

func (l *simpleLogger) Criticalf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, CriticalLevel, format, a...)
}

func (l *simpleLogger) Errorf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, ErrorLevel, format, a...)
}

func (l *simpleLogger) Warnf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, WarnLevel, format, a...)
}

func (l *simpleLogger) Infof(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, InfoLevel, format, a...)
}

func (l *simpleLogger) Debugf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, DebugLevel, format, a...)
}

func (l *simpleLogger) Tracef(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, TraceLevel, format, a...)
}

func (l *simpleLogger) printf(ctx context.Context, lv Level, format string, a ...interface{}) {
	l.Print(ctx, lv, fmt.Sprintf(format, a...), nil)
}

func (l *simpleLogger) Print(ctx context.Context, lv Level, msg string, labels map[string]string) {
	if !shouldPrint(l.level, lv) {
		return
	}

	labelMsg := ""
	if len(labels) > 0 {
		keys := make([]string, len(labels))
		i := 0
		for key := range labels {
			keys[i] = key
			i++
		}
		sort.Strings(keys)

		ls := make([]string, len(labels))
		for i, key := range keys {
			ls[i] = fmt.Sprintf("%s: %s", key, labels[key])
		}
		labelMsg = fmt.Sprintf(" (%s)", strings.Join(ls, ", "))
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, msg+labelMsg)
}
