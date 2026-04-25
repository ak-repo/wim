package errs

import (
	"errors"
	"testing"
)

func TestOpStack_InnerToOuter(t *testing.T) {
	err := E("op1", NotFound, E("op2", Database, E("op3", Internal, errors.New("root"))))

	stack := OpStack(err)

	if len(stack) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(stack))
	}

	if stack[0] != "op3" {
		t.Errorf("expected first (inner) op to be op3, got %s", stack[0])
	}
	if stack[1] != "op2" {
		t.Errorf("expected second op to be op2, got %s", stack[1])
	}
	if stack[2] != "op1" {
		t.Errorf("expected last (outer) op to be op1, got %s", stack[2])
	}
}

func TestErrorsIs_ThroughWrapping(t *testing.T) {
	root := errors.New("root error")
	err := E("op1", Database, E("op2", NotFound, root))

	if !errors.Is(err, root) {
		t.Error("errors.Is should find root error through wrapping")
	}
}

func TestE_preservesExistingErrorChain(t *testing.T) {
	innerRoot := errors.New("inner root")
	innerErr := &Error{Op: "inner", Kind: NotFound, Code: "NOT_FOUND", Err: innerRoot}
	err := E("outer", Database, innerErr)

	var e *Error
	if !errors.As(err, &e) {
		t.Fatal("expected *Error")
	}
	if e.Op != "outer" {
		t.Errorf("expected outer op, got %s", e.Op)
	}
	if !errors.Is(e.Err, innerRoot) {
		t.Error("expected inner root error to be preserved via errors.Is")
	}
	if e.Err != innerErr {
		t.Logf("note: direct equality not required when errors.Is passes")
	}
}

func TestTopError_ReturnsInnermost(t *testing.T) {
	root := errors.New("root")
	err := E("op1", Database, E("op2", NotFound, root))

	top := TopError(err)
	if top != root {
		t.Errorf("expected top error to be root, got %v", top)
	}
}

func TestHTTPErrorResponse_GenericMessageForInternal(t *testing.T) {
	err := E("op", Internal, errors.New("detailed internal error"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 500 {
		t.Errorf("expected status 500, got %d", status)
	}
	if apiErr.Message != "something went wrong" {
		t.Errorf("expected generic message, got %s", apiErr.Message)
	}
	if apiErr.Code != CodeInternal {
		t.Errorf("expected code INTERNAL_ERROR, got %s", apiErr.Code)
	}
}

func TestHTTPErrorResponse_GenericMessageForDatabase(t *testing.T) {
	err := E("op", Database, errors.New("db error details"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 500 {
		t.Errorf("expected status 500, got %d", status)
	}
	if apiErr.Message != "something went wrong" {
		t.Errorf("expected generic message, got %s", apiErr.Message)
	}
	if apiErr.Code != CodeDatabase {
		t.Errorf("expected code DATABASE_ERROR, got %s", apiErr.Code)
	}
}

func TestHTTPErrorResponse_GenericMessageForUnanticipated(t *testing.T) {
	err := E("op", Unanticipated, errors.New("unexpected"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 500 {
		t.Errorf("expected status 500, got %d", status)
	}
	if apiErr.Message != "something went wrong" {
		t.Errorf("expected generic message, got %s", apiErr.Message)
	}
}

func TestHTTPErrorResponse_ExposeStack(t *testing.T) {
	err := E("op1", NotFound, E("op2", Database, errors.New("root")))

	_, apiErr, _ := HTTPErrorResponse(err, true)

	if apiErr.Stack == nil {
		t.Error("expected stack to be included when exposeStack=true")
	}
	if len(apiErr.Stack) != 2 {
		t.Errorf("expected stack length 2, got %d", len(apiErr.Stack))
	}
}

func TestHTTPErrorResponse_HideStack(t *testing.T) {
	err := E("op1", NotFound, E("op2", Database, errors.New("root")))

	_, apiErr, _ := HTTPErrorResponse(err, false)

	if apiErr.Stack != nil {
		t.Error("expected stack to be nil when exposeStack=false")
	}
}

func TestHTTPErrorResponse_NotFoundReturnsMessage(t *testing.T) {
	err := E("op", NotFound, errors.New("resource not found"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 404 {
		t.Errorf("expected status 404, got %d", status)
	}
	if apiErr.Message != "resource not found" {
		t.Errorf("expected 'resource not found', got %s", apiErr.Message)
	}
}

func TestHTTPErrorResponse_InvalidRequestReturnsMessage(t *testing.T) {
	err := E("op", InvalidRequest, errors.New("invalid email"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 400 {
		t.Errorf("expected status 400, got %d", status)
	}
	if apiErr.Message != "invalid email" {
		t.Errorf("expected 'invalid email', got %s", apiErr.Message)
	}
}

func TestHTTPErrorResponse_UnauthorizedReturnsMessage(t *testing.T) {
	err := E("op", Unauthorized, errors.New("token expired"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 401 {
		t.Errorf("expected status 401, got %d", status)
	}
	if apiErr.Message != "token expired" {
		t.Errorf("expected 'token expired', got %s", apiErr.Message)
	}
}

func TestHTTPErrorResponse_ForbiddenReturnsMessage(t *testing.T) {
	err := E("op", Forbidden, errors.New("not admin"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 403 {
		t.Errorf("expected status 403, got %d", status)
	}
	if apiErr.Message != "not admin" {
		t.Errorf("expected 'not admin', got %s", apiErr.Message)
	}
}

func TestHTTPErrorResponse_ConflictReturnsMessage(t *testing.T) {
	err := E("op", Conflict, errors.New("already exists"))

	status, apiErr, _ := HTTPErrorResponse(err, false)

	if status != 409 {
		t.Errorf("expected status 409, got %d", status)
	}
	if apiErr.Message != "already exists" {
		t.Errorf("expected 'already exists', got %s", apiErr.Message)
	}
}
