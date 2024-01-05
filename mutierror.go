package errorx

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
	errors := make([]error, 0, len(e))
	for i := range e {
		errors = append(errors, e[i].Unwrap())
	}
	return errors
}
