package endpoint

import (
	"context"

	"github.com/RangelReale/go-kit-typed/util"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

// Adapter is an adapter from a standard go-kit endpoint to the typed version.
func Adapter[Req any, Resp any](endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return func(ctx context.Context, request Req) (Resp, error) {
		return util.ReturnTypeWithError[Resp](endpoint(ctx, request))
	}
}

// ReverseAdapter is an adapter from a typed endpoint to a standard go-kit version.
func ReverseAdapter[Req any, Resp any](e Endpoint[Req, Resp]) gokitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return util.CallTypeResponseWithError[Req, interface{}](request, func(r Req) (interface{}, error) {
			return e(ctx, r)
		})
	}
}

// Cast casts the standard go-kit endpoint to a typed endpoint of the same type as the first
// parameter.
func Cast[Req any, Resp any](_ Endpoint[Req, Resp], endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return Adapter[Req, Resp](endpoint)
}
