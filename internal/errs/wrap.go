package errs

import "errors"

type Option func(*Error)

func WithCode(code string) Option {
	return func(e *Error) {
		e.Code = code
	}
}

func WithParam(param string) Option {
	return func(e *Error) {
		e.Param = param
	}
}

func E(op Op, kind Kind, err error, opts ...Option) error {
	if err == nil {
		return nil
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		wrapped := &Error{
			Op:    op,
			Kind:  kind,
			Code:  appErr.Code,
			Param: appErr.Param,
			Err:   err,
		}
		for _, opt := range opts {
			opt(wrapped)
		}
		return wrapped
	}

	wrapped := &Error{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
	for _, opt := range opts {
		opt(wrapped)
	}
	return wrapped
}