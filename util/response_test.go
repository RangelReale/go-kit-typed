package util

import (
	"errors"
	"fmt"
	"testing"
)

func TestReturnTypeWithError(t *testing.T) {
	errTest := errors.New("test")

	tests := []struct {
		f        func() (any, error)
		expected int
		error    error
	}{
		{
			f: func() (any, error) {
				return 12, nil
			},
			expected: 12,
		},
		{
			f: func() (any, error) {
				return nil, nil
			},
			expected: 0,
		},
		{
			f: func() (any, error) {
				return "", nil
			},
			error: ErrParameterInvalidType,
		},
		{
			f: func() (any, error) {
				return 0, errTest
			},
			error: errTest,
		},
		{
			f: func() (any, error) {
				return "abc", errTest
			},
			error: errTest,
		},
	}

	for _, test := range tests {
		resp, err := ReturnTypeWithError[int](test.f())
		if test.error != nil {
			if !errors.Is(err, test.error) {
				t.Fatalf("expected '%v' got '%v'", test.error, err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if resp != test.expected {
				t.Errorf("want %d, have %d", test.expected, resp)
			}
		}
	}
}

func TestCallTypeWithError(t *testing.T) {
	errTest := errors.New("test")

	tests := []struct {
		i         interface{}
		funcError error
		expected  int
		error     error
	}{
		{
			i:        12,
			expected: 12,
		},
		{
			i:        nil,
			expected: 0,
		},
		{
			i:     "",
			error: ErrParameterInvalidType,
		},
		{
			i:         10,
			funcError: errTest,
			error:     errTest,
		},
		{
			i:         nil,
			funcError: errTest,
			error:     errTest,
		},
		{
			i:         "",
			funcError: errTest,
			error:     ErrParameterInvalidType,
		},
	}

	for _, test := range tests {
		var resp int
		err := CallTypeWithError[int](test.i, func(r int) error {
			resp = r
			return test.funcError
		})
		if test.error != nil {
			if !errors.Is(err, test.error) {
				t.Fatalf("expected '%v' got '%v'", test.error, err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if resp != test.expected {
				t.Errorf("want %d, have %d", test.expected, resp)
			}
		}
	}
}

func TestCallTypeResponseWithError(t *testing.T) {
	errTest := errors.New("test")

	tests := []struct {
		i         interface{}
		funcError error
		expected  string
		error     error
	}{
		{
			i:        12,
			expected: "12",
		},
		{
			i:        nil,
			expected: "0",
		},
		{
			i:     "",
			error: ErrParameterInvalidType,
		},
		{
			i:         10,
			funcError: errTest,
			error:     errTest,
		},
		{
			i:         nil,
			funcError: errTest,
			error:     errTest,
		},
		{
			i:         "",
			funcError: errTest,
			error:     ErrParameterInvalidType,
		},
	}

	for _, test := range tests {
		resp, err := CallTypeResponseWithError[int, string](test.i, func(r int) (string, error) {
			return fmt.Sprintf("%d", r), test.funcError
		})
		if test.error != nil {
			if !errors.Is(err, test.error) {
				t.Fatalf("expected '%v' got '%v'", test.error, err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if resp != test.expected {
				t.Errorf("want %s, have %s", test.expected, resp)
			}
		}
	}
}
