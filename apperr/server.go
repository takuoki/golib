package apperr

type serverError struct {
	status  int
	code    string
	message string
	log     string
}

// NewServerError creates new server error.
// For RESTful API, set the HTTP status code to `status`.
// Set the error code (ex. "S0001") that the client can handle to `code`.
func NewServerError(status int, code, message, log string) Err {
	return &serverError{
		status:  status,
		code:    code,
		message: message,
		log:     log,
	}
}

// Error is a method to satisfy the error interface.
func (e *serverError) Error() string {
	return e.message
}

// Status returns status value.
func (e *serverError) Status() int {
	return e.status
}

// Code returns code string.
func (e *serverError) Code() string {
	return e.code
}

// Message returns message string.
func (e *serverError) Message() string {
	return e.message
}

// Log returns log string.
func (e *serverError) Log() string {
	return e.log
}

// Type returns error type.
func (e *serverError) Type() Type {
	return ServerError
}
