package util

// ReturnType checks and returns the asked type, or an error
func ReturnType[R any](r any, err error) (R, error) {
	if err != nil {
		var r R
		return r, err
	}
	switch rt := r.(type) {
	case nil:
		var r R
		return r, nil
	case R:
		return rt, nil
	default:
		var r R
		return r, ErrParameterInvalidType
	}
}
