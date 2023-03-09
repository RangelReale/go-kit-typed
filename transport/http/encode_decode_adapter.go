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

// DecodeRequestFuncReverseAdapter is an adapter tp the non-generic DecodeRequestFunc function
func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokithttptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return f(ctx, r)
	}
}

// EncodeRequestFuncReverseAdapter is an adapter to the non-generic EncodeRequestFunc function
func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokithttptransport.EncodeRequestFunc {
	return func(ctx context.Context, request *http.Request, i interface{}) error {
		return util.CallTypeWithError[Req](i, func(r Req) error {
			return f(ctx, request, r)
		})
	}
}

// CreateRequestFuncReverseAdapter is an adapter to the non-generic CreateRequestFunc function
func CreateRequestFuncReverseAdapter[Req any](f CreateRequestFunc[Req]) gokithttptransport.CreateRequestFunc {
	return func(ctx context.Context, i interface{}) (*http.Request, error) {
		var req *http.Request
		err := util.CallTypeWithError[Req](i, func(r Req) error {
			var callErr error
			req, callErr = f(ctx, r)
			return callErr
		})
		return req, err
	}
}

// EncodeResponseFuncReverseAdapter is an adapter to the non-generic EncodeResponseFunc function
func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokithttptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, i interface{}) error {
		return util.CallTypeWithError[Resp](i, func(r Resp) error {
			return f(ctx, w, r)
		})
	}
}

// DecodeResponseFuncReverseAdapter is an adapter to the non-generic DecodeResponseFunc function
func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokithttptransport.DecodeResponseFunc {
	return func(ctx context.Context, r *http.Response) (interface{}, error) {
		return f(ctx, r)
	}
}
