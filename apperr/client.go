package apperr

type clientError struct {
	status  int
	code    string
	message string
}

// NewClientError creates new client error.
// For RESTful API, set the HTTP status code to `status`.
// Set the error code (ex. "E0001") that the client can handle to `code`.
func NewClientError(status int, code, message string) Err {
	return &clientError{
		status:  status,
		code:    code,
		message: message,
	}
}

// Error is a method to satisfy the error interface.
func (e *clientError) Error() string {
	return e.message
}

// Status returns status value.
func (e *clientError) Status() int {
	return e.status
}

// Code returns code string.
func (e *clientError) Code() string {
	return e.code
}

// Message returns message string.
func (e *clientError) Message() string {
	return e.message
}

// Log returns log string.
func (e *clientError) Log() string {
	return ""
}

// Type returns error type.
func (e *clientError) Type() Type {
	return ClientError
}
