package errorx

type CustomError interface {
	error
	With(key string, value interface{}) CustomError
	WithData(map[string]any) CustomError
	Cause() error
	Unwrap() error
}

var _ CustomError = &customError{}
