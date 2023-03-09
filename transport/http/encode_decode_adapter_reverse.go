package http

import (
	"context"
	"net/http"

	"github.com/RangelReale/go-kit-typed/util"
	gokithttptransport "github.com/go-kit/kit/transport/http"
)

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
		return util.CallTypeResponseWithError[Req, *http.Request](i, func(r Req) (*http.Request, error) {
			return f(ctx, r)
		})
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
