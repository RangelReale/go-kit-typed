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
			return rr, util.ErrParameterInvalidType
		}
	}
}

// EndpointReverseAdapter is an adapter from a typed endpoint to a standard go-kit version.
func EndpointReverseAdapter[Req any, Resp any](e Endpoint[Req, Resp]) gokitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		switch tr := request.(type) {
		case nil:
			var r Req
			return e(ctx, r)
		case Req:
			return e(ctx, tr)
		default:
			var r Req
			return r, util.ErrParameterInvalidType
		}
	}
}

// EndpointCast casts the standard go-kit endpoint to a typed endpoint of the same type as the first
// parameter.
func EndpointCast[Req any, Resp any](_ Endpoint[Req, Resp], endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return EndpointAdapter[Req, Resp](endpoint)
}

// MiddlewareWrapper is an adapter for middlewares and generic endpoints.
func MiddlewareWrapper[Req any, Resp any](middleware gokitendpoint.Middleware,
	endpoint Endpoint[Req, Resp]) Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		rmw := middleware(func(ctx context.Context, mreq interface{}) (interface{}, error) {
			switch tr := mreq.(type) {
			case nil:
				var r Req
				return endpoint(ctx, r)
			case Req:
				return endpoint(ctx, tr)
			default:
				var r Req
				return r, util.ErrParameterInvalidType
			}
		})
		resp, err := rmw(ctx, req)
		if err != nil {
			var r Resp
			return r, err
		}
		switch rt := resp.(type) {
		case nil:
			var r Resp
			return r, nil
		case Resp:
			return rt, nil
		default:
			var r Resp
			return r, util.ErrParameterInvalidType
		}
	}
}

// MiddlewareAdapter is an adapter for middlewares and generic endpoints.
func MiddlewareAdapter[Req any, Resp any](middleware gokitendpoint.Middleware) Middleware[Req, Resp] {
	return func(next Endpoint[Req, Resp]) Endpoint[Req, Resp] {
		return func(ctx context.Context, request Req) (Resp, error) {
			resp, err := middleware(EndpointReverseAdapter(next))(ctx, request)
			if err != nil {
				var r Resp
				return r, err
			}
			switch rt := resp.(type) {
			case nil:
				var r Resp
				return r, nil
			case Resp:
				return rt, nil
			default:
				var r Resp
				return r, util.ErrParameterInvalidType
			}
		}
	}
}
