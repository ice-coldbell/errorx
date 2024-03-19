package errorx

func Join(errs ...error) customErrors {
	var cErrs customErrors
	for _, err := range errs {
		switch tErr := err.(type) {
		case nil:
			continue
		case customErrors:
			cErrs = append(cErrs, tErr...)
		case *customError:
			cErrs = append(cErrs, tErr)
		default:
			cErrs = append(cErrs, WrapDepth(err, 4))
		}
	}
	if len(cErrs) == 0 {
		return nil
	}
	return cErrs
}
