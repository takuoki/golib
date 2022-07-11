package apperr_test

import (
	"errors"
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

func ExampleExtractFromGRPCError() {

	callGrpcAPI := func(e apperr.Err) error {
		return e.GRPCError("domain")
	}

	// Client error
	err := callGrpcAPI(NotFound)
	if e, ok := apperr.ExtractFromGRPCError(err); ok {
		fmt.Printf("Code: %s, DetailCode: %s, Message: %s\n", e.Code(), e.DetailCode(), e.Message())
	}

	// Server error
	err = callGrpcAPI(NewInternalServerError(errors.New("error")))
	if e, ok := apperr.ExtractFromGRPCError(err); ok {
		fmt.Printf("Code: %s, DetailCode: %s, Message: %s\n", e.Code(), e.DetailCode(), e.Message())
	}

	// nil error
	err = nil
	if e, ok := apperr.ExtractFromGRPCError(err); ok {
		fmt.Printf("Code: %s, DetailCode: %s, Message: %s\n", e.Code(), e.DetailCode(), e.Message())
	}

	// non-gRPC error
	err = errors.New("error")
	if e, ok := apperr.ExtractFromGRPCError(err); ok {
		fmt.Printf("Code: %s, DetailCode: %s, Message: %s\n", e.Code(), e.DetailCode(), e.Message())
	}

	// Output:
	// Code: NotFound, DetailCode: E0001, Message: not found
	// Code: Internal, DetailCode: S0001, Message: internal server error
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
