package errorx

import "errors"

func Wrap(err error) CustomError {
	return WrapDepth(err, 4)
}

func WrapDepth(err error, depth int) CustomError {
	if err == nil {
		return nil
	}
	if realErr, ok := err.(*customError); ok {
		return realErr
	}
	return &customError{
		err:   err,
		stack: callers(depth),
		data:  make(map[string]any),
	}
}

func WrapWithData(err error, data map[string]any) CustomError {
	return wrapWithData(err, data)
}

func wrapWithData(err error, data map[string]any) CustomError {
	if err == nil {
		return nil
	}
	if realErr, ok := err.(*customError); ok {
		return realErr.WithData(data)
	}
	return &customError{
		err:   err,
		stack: callers(4),
		data:  data,
	}
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
