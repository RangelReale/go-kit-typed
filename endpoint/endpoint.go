package endpoint

import (
	"context"

	gokitendpoint "github.com/go-kit/kit/endpoint"
)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint[Req any, Resp any] func(ctx context.Context, request Req) (response Resp, err error)

func EndpointAdapter[Req any, Resp any](endpoint gokitendpoint.Endpoint) Endpoint[Req, Resp] {
	return func(ctx context.Context, request Req) (Resp, error) {
		resp, err := endpoint(ctx, request)
		if resp != nil {
			return resp.(Resp), err
		}
		var r Resp
		return r, nil
	}
}

// MiddlewareAdapter is an adapter for middlewares and generic endpoints
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
