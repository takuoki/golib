package apperr_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/takuoki/golib/apperr"
)

func ExampleError() {
	// buisiness logic
	err := func() error {
		return NotFound
	}()

	// error handling
	e, ok := apperr.Extract(err)
	if !ok {
		e = ServerError(err)
	}
	if e.Log != "" {
		log.Println(e.Log)
	}
	resp := struct {
		Code    string
		Message string
	}{
		Code:    e.Code,
		Message: e.Message,
	}

	fmt.Printf("%+v", resp)

	// Output:
	// {Code:E0001 Message:not found}
}

// The following is assumed to be defined in each application.

var (
	NotFound   = &apperr.Err{Status: http.StatusNotFound, Code: "E0001", Message: "not found"}
	IDRequired = &apperr.Err{Status: http.StatusBadRequest, Code: "E0002", Message: "id is required"}
)

func ServerError(err error) *apperr.Err {
	return &apperr.Err{
		Status:  http.StatusInternalServerError,
		Code:    "S0001",
		Message: "internal server error",
		Log:     err.Error(),
	}
}
