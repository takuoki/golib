package applog_test

import (
	"context"
	"os"

	"github.com/takuoki/golib/applog"
)

func Example() {

	// Get from environment variables, etc.
	logLevel := "WARN"
	imageTag := "v1.0.0"

	// Create logger
	logger := applog.NewBasicLogger(
		os.Stdout,
		applog.LevelOption(applog.ParseLevel(logLevel)),
		applog.TimeFormatOption("YYYY-MM-DD HH:mm:ss"), // This is invalid format for example test.
		applog.ImageTagOption(imageTag),
	)

	// Print log
	ctx := context.Background()
	logger.Error(ctx, "error message")

	// Output:
	// {"time":"YYYY-MM-DD HH:mm:ss","level":"ERROR","message":"error message","image_tag":"v1.0.0"}
}
