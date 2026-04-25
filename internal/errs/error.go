package errs

type Op string

type Kind int

type Error struct {
	Op    Op
	Kind  Kind
	Code  string
	Param string
	Err   error
}

func (k Kind) String() string {
	switch k {
	case InvalidRequest:
		return "INVALID_REQUEST"
	case Unauthorized:
		return "UNAUTHORIZED"
	case Forbidden:
		return "FORBIDDEN"
	case NotFound:
		return "NOT_FOUND"
	case Conflict:
		return "CONFLICT"
	case Database:
		return "DATABASE"
	case Internal:
		return "INTERNAL"
	case Unanticipated:
		return "UNANTICIPATED"
	default:
		return "UNKNOWN"
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code
}

func (e *Error) Unwrap() error {
	return e.Err
}

const (
	Unknown Kind = iota
	InvalidRequest
	Unauthorized
	Forbidden
	NotFound
	Conflict
	Database
	Internal
	Unanticipated
)
