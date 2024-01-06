package errorx

import (
	"errors"
)

type customErrors []CustomError

func (e customErrors) Error() string {
	var b []byte
	for i, err := range e {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (e customErrors) Unwrap() []error {
	errs := make([]error, 0, len(e))
	for i := range e {
		errs = append(errs, e[i])
	}
	return errs
}

func (e customErrors) Is(target error) bool {
	targetMultiError, ok := target.(customErrors)
	if !ok {
		return false
	}

	if len(e) != len(targetMultiError) {
		return false
	}

	for idx := range e.Unwrap() {
		if !errors.Is(e[idx].Unwrap(), targetMultiError[idx].Unwrap()) {
			return false
		}
	}
	return true
}
