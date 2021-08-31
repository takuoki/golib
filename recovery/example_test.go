package recovery_test

import (
	"fmt"
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/takuoki/golib/recovery"
)

func Example_goroutine() {

	eg := errgroup.Group{}
	for i := 0; i < 1; i++ {
		i := i
		eg.Go(func() (err error) {
			panicked := true
			defer func() {
				if r := recover(); r != nil || panicked {
					err = recovery.Recovery(r)
				}
			}()

			fmt.Printf("your code: %d\n", i)
			panicked = false
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// your code: 0
}

func Example_grpcMiddleware() {

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recovery.Recovery),
	}

	_ = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)
}
