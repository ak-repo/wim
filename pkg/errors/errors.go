package errors

import "errors"

var (
	//common
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternal          = errors.New("internal error")
	ErrConcurrentUpdate  = errors.New("concurrent update detected")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidOperation  = errors.New("invalid operation")
	ErrCheckingFaild     = errors.New("failed to check existing")

)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(err error, code, message string) *Error {
	return &Error{Code: code, Message: message, Err: err}
}
