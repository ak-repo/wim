package errs

import "errors"

func OpStack(err error) []string {
	var ops []string
	for {
		var e *Error
		if !errors.As(err, &e) {
			break
		}
		if e.Op != "" {
			ops = append(ops, string(e.Op))
		}
		if e.Err == nil {
			break
		}
		err = e.Err
	}

	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

func TopError(err error) error {
	last := err
	for {
		var e *Error
		if !errors.As(err, &e) {
			break
		}
		if e.Err == nil {
			break
		}
		last = e.Err
		err = e.Err
	}
	return last
}