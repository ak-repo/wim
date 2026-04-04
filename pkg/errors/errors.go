package errors

import "errors"

// --------------------
// STANDARD ERROR CODES
// --------------------
const (
	CodeNotFound          = "NOT_FOUND"
	CodeAlreadyExists     = "ALREADY_EXISTS"
	CodeInvalidInput      = "INVALID_INPUT"
	CodeUnauthorized      = "UNAUTHORIZED"
	CodeForbidden         = "FORBIDDEN"
	CodeInternal          = "INTERNAL_ERROR"
	CodeConcurrentUpdate  = "CONCURRENT_UPDATE"
	CodeInsufficientStock = "INSUFFICIENT_STOCK"
	CodeInvalidOperation  = "INVALID_OPERATION"
	CodeCheckFailed       = "CHECK_FAILED"
	CodeDatabase          = "DATABASE_ERROR"
	CodeRefCodeFailed     = "REFCODE_GENERATION_FAILED"
)

// --------------------
// BASE SENTINEL ERRORS
// --------------------
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternal          = errors.New("internal error")
	ErrConcurrentUpdate  = errors.New("concurrent update detected")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidOperation  = errors.New("invalid operation")
	ErrCheckingFailed    = errors.New("failed to check existing")
	ErrDatabase          = errors.New("database error")

	// specific
	ErrRefCodeGeneration = errors.New("failed to generate refcode")
)

// --------------------
// STRUCTURED ERROR
// --------------------
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // internal only
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

//
// --------------------
// HELPERS
// --------------------
//

// New creates a business error
func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps internal error with safe message
func Wrap(err error, code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

//
// --------------------
// TYPE CHECKING
// --------------------
//

// GetError extracts custom error
func GetError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}
