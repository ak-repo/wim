package errs

import (
	"errors"
	"net/http"

	legacyerrs "github.com/ak-repo/wim/pkg/errors"
)

type ApiError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Stack   []string `json:"stack,omitempty"`
}

func HTTPErrorResponse(err error, exposeStack bool) (int, ApiError, *Error) {
	if err == nil {
		return http.StatusInternalServerError, ApiError{Code: CodeInternal, Message: "something went wrong"}, nil
	}

	e := normalizeError(err)

	status := statusFromKind(e.Kind)
	message := e.messageForKind()
	code := e.Code
	if code == "" {
		code = codeFromKind(e.Kind)
	}

	stack := OpStack(err)
	if !exposeStack {
		stack = nil
	}

	return status, ApiError{Code: code, Message: message, Stack: stack}, e
}

func (e *Error) messageForKind() string {
	switch e.Kind {
	case Internal, Database, Unanticipated:
		return "something went wrong"
	default:
		if e.Err == nil {
			return "unknown error"
		}
		return e.Err.Error()
	}
}

func normalizeError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}

	var le *legacyerrs.Error
	if errors.As(err, &le) {
		return &Error{Kind: legacyCodeToKind(le.Code), Code: le.Code, Err: le}
	}

	if errors.Is(err, legacyerrs.ErrUnauthorized) {
		return &Error{Kind: Unauthorized, Code: CodeUnauthorized, Err: err}
	}
	if errors.Is(err, legacyerrs.ErrForbidden) {
		return &Error{Kind: Forbidden, Code: CodeForbidden, Err: err}
	}
	if errors.Is(err, legacyerrs.ErrInvalidInput) {
		return &Error{Kind: InvalidRequest, Code: CodeInvalidRequest, Err: err}
	}
	if errors.Is(err, legacyerrs.ErrNotFound) {
		return &Error{Kind: NotFound, Code: CodeNotFound, Err: err}
	}
	if errors.Is(err, legacyerrs.ErrDatabase) {
		return &Error{Kind: Database, Code: CodeDatabase, Err: err}
	}

	return &Error{Err: err, Kind: Internal}
}

func legacyCodeToKind(code string) Kind {
	switch code {
	case legacyerrs.CodeInvalidInput:
		return InvalidRequest
	case legacyerrs.CodeUnauthorized:
		return Unauthorized
	case legacyerrs.CodeForbidden:
		return Forbidden
	case legacyerrs.CodeNotFound:
		return NotFound
	case legacyerrs.CodeAlreadyExists, legacyerrs.CodeConcurrentUpdate:
		return Conflict
	case legacyerrs.CodeDatabase:
		return Database
	default:
		return Internal
	}
}

func statusFromKind(kind Kind) int {
	switch kind {
	case InvalidRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	case Database, Internal, Unanticipated:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func codeFromKind(kind Kind) string {
	switch kind {
	case InvalidRequest:
		return CodeInvalidRequest
	case Unauthorized:
		return CodeUnauthorized
	case Forbidden:
		return CodeForbidden
	case NotFound:
		return CodeNotFound
	case Conflict:
		return CodeConflict
	case Database:
		return CodeDatabase
	case Internal, Unanticipated:
		return CodeInternal
	default:
		return CodeInternal
	}
}
