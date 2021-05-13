package apperr

// Type is error type.
type Type int

// Type list.
const (
	ClientError Type = iota + 1
	ServerError
)

// String returns a string of the error type.
func (t Type) String() string {
	switch t {
	case ClientError:
		return "client-error"
	case ServerError:
		return "server-error"
	default:
		return "unknown"
	}
}
