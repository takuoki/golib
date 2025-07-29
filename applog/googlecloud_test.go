package applog

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/takuoki/golib/appctx"
)

func TestNewGoogleCloudLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewGoogleCloudLogger(&buf)

	if logger == nil {
		t.Error("NewGoogleCloudLogger should return a logger instance")
	}
}

func TestGoogleCloudLogger_Print(t *testing.T) {
	var buf bytes.Buffer
	logger := NewGoogleCloudLogger(&buf, LevelOption(DebugLevel))

	ctx := context.Background()
	ctx = appctx.WithRequestID(ctx, "test-request-id")

	logger.Info(ctx, "test message")

	var logEntry googleCloudLog
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if logEntry.Severity != "INFO" {
		t.Errorf("Expected severity INFO, got %s", logEntry.Severity)
	}

	if logEntry.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", logEntry.Message)
	}

	if logEntry.Labels["request_id"] != "test-request-id" {
		t.Errorf("Expected request_id in labels to be 'test-request-id', got %s", logEntry.Labels["request_id"])
	}
}

func TestGoogleCloudLogger_WithImageTag(t *testing.T) {
	var buf bytes.Buffer
	logger := NewGoogleCloudLogger(&buf, ImageTagOption("v1.2.3"))

	ctx := context.Background()
	logger.Info(ctx, "test message")

	var logEntry googleCloudLog
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if logEntry.Labels["image_tag"] != "v1.2.3" {
		t.Errorf("Expected image_tag in labels to be 'v1.2.3', got %s", logEntry.Labels["image_tag"])
	}
}

func TestGoogleCloudLogger_WithCustomLabels(t *testing.T) {
	var buf bytes.Buffer
	logger := NewGoogleCloudLogger(&buf)

	ctx := context.Background()
	customLabels := map[string]string{
		"service": "api",
		"version": "1.0.0",
	}

	logger.Print(ctx, InfoLevel, "test message", customLabels)

	var logEntry googleCloudLog
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if logEntry.Labels["service"] != "api" {
		t.Errorf("Expected service in labels to be 'api', got %s", logEntry.Labels["service"])
	}

	if logEntry.Labels["version"] != "1.0.0" {
		t.Errorf("Expected version in labels to be '1.0.0', got %s", logEntry.Labels["version"])
	}
}

func TestGoogleCloudLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewGoogleCloudLogger(&buf, LevelOption(WarnLevel))

	ctx := context.Background()

	// Should not be printed (below warn level)
	logger.Info(ctx, "info message")
	logger.Debug(ctx, "debug message")

	// Should be printed (warn level or above)
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, "error message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have 2 lines (warn and error)
	if len(lines) != 2 {
		t.Errorf("Expected 2 log lines, got %d", len(lines))
	}

	// Check that warn message is present
	if !strings.Contains(output, "warn message") {
		t.Error("Expected warn message to be logged")
	}

	// Check that error message is present
	if !strings.Contains(output, "error message") {
		t.Error("Expected error message to be logged")
	}

	// Check that info and debug messages are not present
	if strings.Contains(output, "info message") {
		t.Error("Info message should not be logged when level is WARN")
	}

	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not be logged when level is WARN")
	}
}

func TestGoogleCloudLogger_Formatted(t *testing.T) {
	var buf bytes.Buffer
	logger := NewGoogleCloudLogger(&buf)

	ctx := context.Background()
	logger.Infof(ctx, "formatted message: %s, number: %d", "test", 42)

	var logEntry googleCloudLog
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	expected := "formatted message: test, number: 42"
	if logEntry.Message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, logEntry.Message)
	}
}

func TestLevelToGoogleCloudSeverity(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{CriticalLevel, "CRITICAL"},
		{ErrorLevel, "ERROR"},
		{WarnLevel, "WARNING"},
		{InfoLevel, "INFO"},
		{DebugLevel, "DEBUG"},
		{TraceLevel, "DEBUG"},
		{UnknownLevel, "DEFAULT"},
	}

	for _, test := range tests {
		result := levelToGoogleCloudSeverity(test.level)
		if result != test.expected {
			t.Errorf("For level %v, expected %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestGoogleCloudLogger_TimeFormat(t *testing.T) {
	var buf bytes.Buffer
	customTimeFormat := "2006-01-02T15:04:05Z"
	logger := NewGoogleCloudLogger(&buf, TimeFormatOption(customTimeFormat))

	ctx := context.Background()
	logger.Info(ctx, "test message")

	var logEntry googleCloudLog
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	// Check if timestamp can be parsed with the custom format
	if _, err := time.Parse(customTimeFormat, logEntry.Timestamp); err != nil {
		t.Errorf("Timestamp should be in custom format %s, but got parsing error: %v", customTimeFormat, err)
	}
}
