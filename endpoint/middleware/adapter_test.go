package middleware_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/endpoint/middleware"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

func TestAdapter(t *testing.T) {
	buf := strings.Builder{}

	var e endpoint.Endpoint[string, string]
	{
		e = strendpoint(&buf)
		e = middleware.Adapter[string, string](strmiddleware("third", &buf))(e)
		e = middleware.Adapter[string, string](strmiddleware("second", &buf))(e)
		e = middleware.Adapter[string, string](strmiddleware("first", &buf))(e)
	}

	resp, err := e(ctx, "data")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp != "endpoint-response" {
		t.Errorf("want %s, have %s", "endpoint-response", resp)
	}

	expected := "|pre-first|pre-second|pre-third|endpoint-data|post-third|post-second|post-first"
	if buf.String() != expected {
		t.Errorf("want '%s', have '%s'", expected, buf.String())
	}
}

func TestAdapterChain(t *testing.T) {
	buf := strings.Builder{}

	m := gokitendpoint.Chain(
		strmiddleware("first", &buf),
		strmiddleware("second", &buf),
		strmiddleware("third", &buf),
	)

	e := middleware.Adapter[string, string](m)(strendpoint(&buf))

	resp, err := e(ctx, "data")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp != "endpoint-response" {
		t.Errorf("want %s, have %s", "endpoint-response", resp)
	}

	expected := "|pre-first|pre-second|pre-third|endpoint-data|post-third|post-second|post-first"
	if buf.String() != expected {
		t.Errorf("want '%s', have '%s'", expected, buf.String())
	}
}

func TestWrapper(t *testing.T) {
	buf := strings.Builder{}

	var e endpoint.Endpoint[string, string]
	{
		e = strendpoint(&buf)
		e = middleware.Wrapper(strmiddleware("third", &buf), e)
		e = middleware.Wrapper(strmiddleware("second", &buf), e)
		e = middleware.Wrapper(strmiddleware("first", &buf), e)
	}

	resp, err := e(ctx, "data")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp != "endpoint-response" {
		t.Errorf("want %s, have %s", "endpoint-response", resp)
	}

	expected := "|pre-first|pre-second|pre-third|endpoint-data|post-third|post-second|post-first"
	if buf.String() != expected {
		t.Errorf("want '%s', have '%s'", expected, buf.String())
	}
}

func TestWrapperChain(t *testing.T) {
	buf := strings.Builder{}

	m := gokitendpoint.Chain(
		strmiddleware("first", &buf),
		strmiddleware("second", &buf),
		strmiddleware("third", &buf),
	)

	e := middleware.Wrapper(m, strendpoint(&buf))

	resp, err := e(ctx, "data")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp != "endpoint-response" {
		t.Errorf("want %s, have %s", "endpoint-response", resp)
	}

	expected := "|pre-first|pre-second|pre-third|endpoint-data|post-third|post-second|post-first"
	if buf.String() != expected {
		t.Errorf("want '%s', have '%s'", expected, buf.String())
	}
}

func strmiddleware(s string, buf *strings.Builder) gokitendpoint.Middleware {
	return func(next gokitendpoint.Endpoint) gokitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			buf.WriteString(fmt.Sprintf("|pre-%s", s))
			ret, err := next(ctx, request)
			buf.WriteString(fmt.Sprintf("|post-%s", s))
			return ret, err
		}
	}
}

func strendpoint(buf *strings.Builder) endpoint.Endpoint[string, string] {
	return func(ctx context.Context, request string) (string, error) {
		buf.WriteString(fmt.Sprintf("|endpoint-%s", request))
		return "endpoint-response", nil
	}
}
