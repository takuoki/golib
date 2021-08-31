// Package recovery defines a function that converts a panic into an error.
package recovery

import (
	"errors"
	"fmt"
	"runtime"
)

// Recovery converts a panic into an error.
// See: example_test.go
func Recovery(p interface{}) error {
	e := fmt.Sprintf("panic recovered: %v", p)
	for depth := 0; ; depth++ {
		pc, src, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		e += fmt.Sprintf("%d: %s: %s(%d)\n", depth, runtime.FuncForPC(pc).Name(), src, line)
	}
	return errors.New(e)
}
