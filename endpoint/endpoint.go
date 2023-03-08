package endpoint

import (
	"context"

	gokitendpoint "github.com/go-kit/kit/endpoint"
)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint[Req any, Resp any] func(ctx context.Context, request Req) (response Resp, err error)

// EndpointAdapter is an adapter from a standard go-kit endpoint to the typed version.
func EndpointAdapter[Req any, Resp any](endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return func(ctx context.Context, request Req) (Resp, error) {
		resp, err := endpoint(ctx, request)
		if err != nil {
			var r Resp
			return r, err
		}

		switch tr := resp.(type) {
		case nil:
			var rr Resp
			return rr, nil
		case Resp:
			return tr, nil
		default:
			var rr Resp
			return rr, ErrParameterInvalidType
		}
	}
}

// EndpointAdapterBack is an adapter from a typed endpoint to a standard go-kit version.
func EndpointAdapterBack[Req any, Resp any](e Endpoint[Req, Resp]) gokitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		switch tr := request.(type) {
		case nil:
			var r Req
			return e(ctx, r)
		case Req:
			return e(ctx, tr)
		default:
			var r Req
			return r, ErrParameterInvalidType
		}
	}
}

// MiddlewareAdapter is an adapter for middlewares and generic endpoints.
func MiddlewareAdapter[Req any, Resp any](middleware gokitendpoint.Middleware,
	endpoint Endpoint[Req, Resp]) Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		rmw := middleware(func(ctx context.Context, mreq interface{}) (interface{}, error) {
			return endpoint(ctx, mreq.(Req))
		})
		resp, err := rmw(ctx, req)
		if resp != nil {
			return resp.(Resp), err
		}
		var r Resp
		return r, nil
	}
}
