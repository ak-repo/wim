package errs

import (
	"errors"
	"fmt"
)

const (
	CodeInvalidInput      = CodeInvalidRequest
	CodeConcurrentUpdate  = "CONCURRENT_UPDATE"
	CodeInsufficientStock = "INSUFFICIENT_STOCK"
	CodeInvalidOperation  = "INVALID_OPERATION"
	CodeCheckFailed       = "CHECK_FAILED"
	CodeRefCodeFailed     = "REFCODE_GENERATION_FAILED"
)

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
	ErrRefCodeGeneration = errors.New("failed to generate refcode")
)

func New(code, message string) error {
	return &Error{
		Kind: codeToKind(code),
		Code: code,
		Err:  errors.New(message),
	}
}

func Wrap(err error, code, message string) error {
	if err == nil {
		return nil
	}
	return &Error{
		Kind: codeToKind(code),
		Code: code,
		Err:  fmt.Errorf("%s: %w", message, err),
	}
}

func codeToKind(code string) Kind {
	switch code {
	case CodeInvalidRequest:
		return InvalidRequest
	case CodeUnauthorized:
		return Unauthorized
	case CodeForbidden:
		return Forbidden
	case CodeNotFound:
		return NotFound
	case CodeAlreadyExists, CodeConflict, CodeConcurrentUpdate:
		return Conflict
	case CodeDatabase:
		return Database
	default:
		return Internal
	}
}
