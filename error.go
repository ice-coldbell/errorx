package errorx

import "errors"

type customError struct {
	err   error
	stack stack
	data  map[string]any
}

func New(message string) CustomError {
	return &customError{
		err:   errors.New(message),
		stack: callers(3),
		data:  make(map[string]any),
	}
}

func (e *customError) Error() string {
	if e == nil {
		return ""
	}
	return e.err.Error()
}

func (e *customError) Cause() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *customError) With(key string, data any) CustomError {
	e.data[key] = data
	return e
}

func (e *customError) WithData(data map[string]any) CustomError {
	for k, v := range data {
		e.data[k] = v
	}
	return e
}

func (e *customError) StackTrace() stack {
	if e == nil {
		return nil
	}
	return e.stack
}

func (e *customError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}
