package endpoint

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/RangelReale/go-kit-typed/util"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

func TestAdapter(t *testing.T) {
	errTest := errors.New("test")

	tests := []struct {
		f        gokitendpoint.Endpoint
		expected string
		error    error
	}{
		{
			f: func(ctx context.Context, request interface{}) (interface{}, error) {
				return fmt.Sprintf("str-%v", request), nil
			},
			expected: "str-data",
		},
		{
			f: func(ctx context.Context, request interface{}) (interface{}, error) {
				return 12, nil
			},
			error: util.ErrParameterInvalidType,
		},
		{
			f: func(ctx context.Context, request interface{}) (interface{}, error) {
				return nil, nil
			},
			expected: "",
		},
		{
			f: func(ctx context.Context, request interface{}) (interface{}, error) {
				return nil, errTest
			},
			error: errTest,
		},
	}

	for _, test := range tests {
		resp, err := Adapter[string, string](test.f)(context.Background(), "data")
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

func TestReverseAdapter(t *testing.T) {
	errTest := errors.New("test")

	tests := []struct {
		f        Endpoint[string, string]
		expected string
		error    error
	}{
		{
			f: func(ctx context.Context, request string) (string, error) {
				return fmt.Sprintf("str-%v", request), nil
			},
			expected: "str-data",
		},
		{
			f: func(ctx context.Context, request string) (string, error) {
				return "", errTest
			},
			error: errTest,
		},
	}

	for _, test := range tests {
		resp, err := ReverseAdapter[string, string](test.f)(context.Background(), "data")
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
