package applog_test

import (
	"context"
	"log"
	"os"

	"github.com/takuoki/golib/appctx"
	"github.com/takuoki/golib/applog"
)

func Example() {

	// Get from environment variables, etc.
	logLevel := "WARN"
	imageTag := "v1.0.0"

	// Create logger
	lv, err := applog.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Fail to parse log level (logLevel=%q): %v", logLevel, err)
	}
	logger := applog.NewBasicLogger(
		os.Stdout,
		applog.LevelOption(lv),
		applog.TimeFormatOption("YYYY-MM-DD HH:mm:ss"), // This is invalid format for example test.
		applog.ImageTagOption(imageTag),
	)

	// Print log
	ctx := context.Background()
	logger.Error(ctx, "error message")

	// Output:
	// {"time":"YYYY-MM-DD HH:mm:ss","level":"ERROR","message":"error message","image_tag":"v1.0.0"}
}

func ExampleNewGoogleCloudLogger() {

	// Get from environment variables, etc.
	logLevel := "WARN"
	imageTag := "v1.0.0"

	// Create Google Cloud Logger
	lv, err := applog.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("Fail to parse log level (logLevel=%q): %v", logLevel, err)
	}
	logger := applog.NewGoogleCloudLogger(
		os.Stdout,
		applog.LevelOption(lv),
		applog.TimeFormatOption("YYYY-MM-DDTHH:mm:ssZ"), // This is invalid format for example test.
		applog.ImageTagOption(imageTag),
	)

	// Print log with context including request ID
	ctx := context.Background()
	ctx = appctx.WithRequestID(ctx, "req-12345")

	// Log with custom labels
	customLabels := map[string]string{
		"service": "api-server",
		"version": "1.2.3",
	}
	logger.Print(ctx, applog.ErrorLevel, "error message", customLabels)

	// Output:
	// {"timestamp":"YYYY-MM-DDTHH:mm:ssZ","severity":"ERROR","message":"error message","labels":{"image_tag":"v1.0.0","request_id":"req-12345","service":"api-server","version":"1.2.3"}}
}
