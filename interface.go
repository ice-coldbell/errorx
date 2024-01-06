package errorx

type CustomError interface {
	error
	With(key string, value interface{}) CustomError
	WithData(map[string]any) CustomError
	Cause() error
	Unwrap() error
}

type CustomErrors interface {
	error
	Is(error) bool
	Unwrap() []error
}

var _ CustomError = &customError{}

var _ CustomErrors = customErrors{}
