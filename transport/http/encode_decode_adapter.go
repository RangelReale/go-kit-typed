package http

import (
	"context"
	"net/http"

	"github.com/RangelReale/go-kit-typed/util"
	gokithttptransport "github.com/go-kit/kit/transport/http"
)

// DecodeRequestFuncAdapter is an adapter from the non-generic DecodeRequestFunc function
func DecodeRequestFuncAdapter[Req any](f gokithttptransport.DecodeRequestFunc) DecodeRequestFunc[Req] {
	return func(ctx context.Context, r *http.Request) (Req, error) {
		return util.ReturnTypeWithError[Req](f(ctx, r))
	}
}

// EncodeRequestFuncAdapter is an adapter from the non-generic EncodeRequestFunc function
func EncodeRequestFuncAdapter[Req any](f gokithttptransport.EncodeRequestFunc) EncodeRequestFunc[Req] {
	return func(ctx context.Context, request *http.Request, req Req) error {
		return f(ctx, request, req)
	}
}

// EncodeResponseFuncAdapter is an adapter from the non-generic EncodeResponseFunc function
func EncodeResponseFuncAdapter[Resp any](f gokithttptransport.EncodeResponseFunc) EncodeResponseFunc[Resp] {
	return func(ctx context.Context, writer http.ResponseWriter, resp Resp) error {
		return f(ctx, writer, resp)
	}
}

// DecodeResponseFuncAdapter is an adapter from the non-generic DecodeRequestFunc function
func DecodeResponseFuncAdapter[Resp any](f gokithttptransport.DecodeResponseFunc) DecodeResponseFunc[Resp] {
	return func(ctx context.Context, r *http.Response) (Resp, error) {
		return util.ReturnTypeWithError[Resp](f(ctx, r))
	}
}
