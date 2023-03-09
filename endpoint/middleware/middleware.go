package middleware

import (
	"context"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

// Adapter is an adapter for middlewares and generic endpoints.
func Adapter[Req any, Resp any](middleware gokitendpoint.Middleware) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(ctx context.Context, request Req) (Resp, error) {
			return util.ReturnTypeWithError[Resp](middleware(endpoint.ReverseAdapter(next))(ctx, request))
		}
	}
}

// Wrapper is a wrapper for middlewares and generic endpoints.
func Wrapper[Req any, Resp any](middleware gokitendpoint.Middleware,
	endpoint endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		rmw := middleware(func(ctx context.Context, mreq interface{}) (interface{}, error) {
			return util.CallTypeResponseWithError[Req, interface{}](mreq, func(r Req) (interface{}, error) {
				return endpoint(ctx, r)
			})
		})
		return util.ReturnTypeWithError[Resp](rmw(ctx, req))
	}
}
