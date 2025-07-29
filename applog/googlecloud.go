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

type googleCloudLogger struct {
	mu         sync.Mutex
	out        io.Writer
	level      Level
	timeFormat string
	imageTag   string
}

// NewGoogleCloudLogger creates a Google Cloud Logging compatible logger that outputs in JSON format.
// It handles the necessity of output according to the log level,
// and outputs context information to the log in common.
// The output format follows Google Cloud Logging structure with severity and labels.
func NewGoogleCloudLogger(w io.Writer, opts ...Option) Logger {
	logger := &googleCloudLogger{
		out:        w,
		timeFormat: time.RFC3339Nano, // Google Cloud Logging prefers RFC3339Nano
	}
	for _, opt := range opts {
		// googleCloudLogger option never returns an error
		_ = opt(logger)
	}
	return logger
}

func (l *googleCloudLogger) setLevel(lv Level) error {
	l.level = lv
	return nil
}

func (l *googleCloudLogger) setTimeFormat(format string) error {
	l.timeFormat = format
	return nil
}

func (l *googleCloudLogger) setImageTag(tag string) error {
	l.imageTag = tag
	return nil
}

func (l *googleCloudLogger) Critical(ctx context.Context, msg string) {
	l.Print(ctx, CriticalLevel, msg, nil)
}

func (l *googleCloudLogger) Error(ctx context.Context, msg string) {
	l.Print(ctx, ErrorLevel, msg, nil)
}

func (l *googleCloudLogger) Warn(ctx context.Context, msg string) {
	l.Print(ctx, WarnLevel, msg, nil)
}

func (l *googleCloudLogger) Info(ctx context.Context, msg string) {
	l.Print(ctx, InfoLevel, msg, nil)
}

func (l *googleCloudLogger) Debug(ctx context.Context, msg string) {
	l.Print(ctx, DebugLevel, msg, nil)
}

func (l *googleCloudLogger) Trace(ctx context.Context, msg string) {
	l.Print(ctx, TraceLevel, msg, nil)
}

func (l *googleCloudLogger) Criticalf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, CriticalLevel, format, a...)
}

func (l *googleCloudLogger) Errorf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, ErrorLevel, format, a...)
}

func (l *googleCloudLogger) Warnf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, WarnLevel, format, a...)
}

func (l *googleCloudLogger) Infof(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, InfoLevel, format, a...)
}

func (l *googleCloudLogger) Debugf(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, DebugLevel, format, a...)
}

func (l *googleCloudLogger) Tracef(ctx context.Context, format string, a ...interface{}) {
	l.printf(ctx, TraceLevel, format, a...)
}

func (l *googleCloudLogger) printf(ctx context.Context, lv Level, format string, a ...interface{}) {
	l.Print(ctx, lv, fmt.Sprintf(format, a...), nil)
}

func (l *googleCloudLogger) Print(ctx context.Context, lv Level, msg string, labels map[string]string) {
	if !shouldPrint(l.level, lv) {
		return
	}

	// Prepare labels map
	logLabels := make(map[string]string)

	// Add imageTag to labels if set
	if l.imageTag != "" {
		logLabels["image_tag"] = l.imageTag
	}

	// Add requestID to labels if available
	if requestID := appctx.RequestID(ctx); requestID != "" {
		logLabels["request_id"] = requestID
	}

	// Add user-provided labels
	for k, v := range labels {
		logLabels[k] = v
	}

	log := googleCloudLog{
		Timestamp: time.Now().Format(l.timeFormat),
		Severity:  levelToGoogleCloudSeverity(lv),
		Message:   msg,
		Labels:    logLabels,
	}

	// Remove labels if empty to keep output clean
	if len(log.Labels) == 0 {
		log.Labels = nil
	}

	jsonLog, _ := json.Marshal(log)

	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.out, string(jsonLog)) //nolint:errcheck
}

// levelToGoogleCloudSeverity converts internal log level to Google Cloud Logging severity
func levelToGoogleCloudSeverity(lv Level) string {
	switch lv {
	case CriticalLevel:
		return "CRITICAL"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARNING"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "DEBUG" // Google Cloud Logging doesn't have TRACE, use DEBUG
	default:
		return "DEFAULT"
	}
}

type googleCloudLog struct {
	Timestamp string            `json:"timestamp"`
	Severity  string            `json:"severity"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels,omitempty"`
}
