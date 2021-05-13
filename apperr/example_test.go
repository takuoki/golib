package apperr_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/takuoki/golib/apperr"
)

func Example() {

	// Buisiness logic
	err := func() error {
		return NotFound
	}()

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
		Code:    e.Code(),
		Message: e.Message(),
	}

	fmt.Printf("%+v", resp)

	// Output:
	// {Code:E0001 Message:not found}
}

// The following is assumed to be defined in each application.

var (
	NotFound   = apperr.NewClientError(http.StatusNotFound, "E0001", "not found")
	IDRequired = apperr.NewClientError(http.StatusBadRequest, "E0002", "id is required")
)

func NewInternalServerError(err error) apperr.Err {
	return apperr.NewServerError(
		http.StatusInternalServerError,
		"S0001",
		"internal server error",
		err.Error(),
	)
}
