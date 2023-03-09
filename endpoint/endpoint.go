package endpoint

import (
	"context"

	"github.com/RangelReale/go-kit-typed/util"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint[Req any, Resp any] func(ctx context.Context, request Req) (response Resp, err error)

// Middleware is a chainable behavior modifier for endpoints.
type Middleware[Req any, Resp any] func(Endpoint[Req, Resp]) Endpoint[Req, Resp]

// Adapter is an adapter from a standard go-kit endpoint to the typed version.
func Adapter[Req any, Resp any](endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return func(ctx context.Context, request Req) (Resp, error) {
		return util.ReturnTypeWithError[Resp](endpoint(ctx, request))
	}
}

// ReverseAdapter is an adapter from a typed endpoint to a standard go-kit version.
func ReverseAdapter[Req any, Resp any](e Endpoint[Req, Resp]) gokitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var resp interface{}
		err := util.CallTypeWithError[Req](request, func(r Req) error {
			var callErr error
			resp, callErr = e(ctx, r)
			return callErr
		})
		return resp, err
	}
}

// Cast casts the standard go-kit endpoint to a typed endpoint of the same type as the first
// parameter.
func Cast[Req any, Resp any](_ Endpoint[Req, Resp], endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return Adapter[Req, Resp](endpoint)
}
