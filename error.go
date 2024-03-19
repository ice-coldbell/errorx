package errorx

import (
	"errors"
	"log/slog"
)

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
	if e == nil || e.err == nil {
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
	if e == nil || e.err == nil {
		return nil
	}
	e.data[key] = data
	return e
}

func (e *customError) WithData(data map[string]any) CustomError {
	if e == nil || e.err == nil {
		return nil
	}
	for k, v := range data {
		e.data[k] = v
	}
	return e
}

// For sentry-go extract stacktrace
func (e *customError) StackTrace() []uintptr {
	if e == nil || e.err == nil {
		return nil
	}
	return e.stack.StackTrace()
}

func (e *customError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *customError) LogValue() slog.Value {
	attr := []slog.Attr{
		slog.String("message", e.Error()),
	}
	for k, v := range e.data {
		attr = append(attr, slog.Any(k, v))
	}
	return slog.GroupValue(attr...)
}
