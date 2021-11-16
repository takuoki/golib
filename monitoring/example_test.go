package monitoring_test

import (
	"context"
	"os"

	"github.com/takuoki/golib/monitoring"
)

func Example() {

	ctx := context.Background()
	ctx = monitoring.Record(ctx, "main start")
	defer func() {
		ctx = monitoring.Record(ctx, "main end")
		monitoring.Output(ctx, os.Stdout, nil)
	}()

	foo(ctx)

	// Output:
	// |    NAME    | ALLOC (MB) | TOTALALLOC (MB) | SYS (MB) |
	// |------------|------------|-----------------|----------|
	// | main start |          0 |               0 |        9 |
	// | foo start  |          0 |               0 |        9 |
	// | foo end    |          0 |               0 |        9 |
	// | main end   |          0 |               0 |        9 |
}

func foo(ctx context.Context) {
	ctx = monitoring.Record(ctx, "foo start")
	defer monitoring.Record(ctx, "foo end")

	// ...
}
