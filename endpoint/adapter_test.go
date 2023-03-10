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
	tests := []struct {
		f           gokitendpoint.Endpoint
		expected    string
		isTypeError bool
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
			isTypeError: true,
		},
		{
			f: func(ctx context.Context, request interface{}) (interface{}, error) {
				return nil, nil
			},
			expected: "",
		},
	}

	for _, test := range tests {
		resp, err := Adapter[string, string](test.f)(context.Background(), "data")
		if test.isTypeError {
			if !errors.Is(err, util.ErrParameterInvalidType) {
				t.Fatalf("expected ErrParameterInvalidType got %v", err)
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
	tests := []struct {
		f        Endpoint[string, string]
		expected string
	}{
		{
			f: func(ctx context.Context, request string) (string, error) {
				return fmt.Sprintf("str-%v", request), nil
			},
			expected: "str-data",
		},
	}

	for _, test := range tests {
		resp, err := ReverseAdapter[string, string](test.f)(context.Background(), "data")
		if err != nil {
			t.Fatal(err)
		}
		if resp != test.expected {
			t.Errorf("want %s, have %s", test.expected, resp)
		}
	}
}
