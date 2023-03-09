package util

// ReturnTypeWithError checks and returns the asked type, or an error
func ReturnTypeWithError[R any](r any, err error) (R, error) {
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

// CallTypeWithError checks and returns the asked type, or an error
func CallTypeWithError[R any](i interface{}, f func(r R) error) error {
	switch ri := i.(type) {
	case nil:
		var rr R
		return f(rr)
	case R:
		return f(ri)
	default:
		return ErrParameterInvalidType
	}
}

// CallTypeResponseWithError checks and returns the asked type, or an error
func CallTypeResponseWithError[R any, Response any](i interface{}, f func(r R) (Response, error)) (Response, error) {
	switch ri := i.(type) {
	case nil:
		var rr R
		return f(rr)
	case R:
		return f(ri)
	default:
		var resp Response
		return resp, ErrParameterInvalidType
	}
}
