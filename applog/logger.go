// Package applog provides a logger intended for use in applications such as APIs.
// We do not require versatility in the output format,
// and we aim to make it simple to use.
// When customizing the output format, do not use this library as it is,
// but refer to this library and create it for each application.
package applog

import "context"

// Logger represents a logging interface that outputs log with log level.
type Logger interface {
	Critical(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
	Debug(ctx context.Context, msg string)
	Trace(ctx context.Context, msg string)

	Criticalf(ctx context.Context, format string, a ...interface{})
	Errorf(ctx context.Context, format string, a ...interface{})
	Warnf(ctx context.Context, format string, a ...interface{})
	Infof(ctx context.Context, format string, a ...interface{})
	Debugf(ctx context.Context, format string, a ...interface{})
	Tracef(ctx context.Context, format string, a ...interface{})

	Print(ctx context.Context, lv Level, msg string, labels map[string]string)

	// Option setter
	SetLevel(lv Level) error
	SetTimeFormat(format string) error
	SetImageTag(tag string) error
}
