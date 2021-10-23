package apperr_test

import (
	"fmt"
	"log"

	"github.com/takuoki/golib/apperr"
	"google.golang.org/grpc/codes"
)

func Example() {

	// Buisiness logic
	err := func(id string) error {
		if id == "" {
			return IDRequired
		}
		return NotFound
	}("id1")

	// Error handling
	e, ok := apperr.Extract(err)
	if !ok {
		e = NewInternalServerError(err)
	}
	if e.Log() != "" {
		log.Println(e.Log())
	}
	resp := struct {
		Code    string
		Message string
	}{
		Code:    e.DetailCode(),
		Message: e.Message(),
	}

	fmt.Printf("%+v", resp)

	// Output:
	// {Code:E0001 Message:not found}
}

// The following is assumed to be defined in each application.

var (
	NotFound   = apperr.NewClientError(codes.NotFound, "E0001", "not found")
	IDRequired = apperr.NewClientError(codes.InvalidArgument, "E0002", "id is required")
)

func NewInternalServerError(err error) apperr.Err {
	return apperr.NewServerError(
		codes.Internal,
		"S0001",
		"internal server error",
		err.Error(),
	)
}
