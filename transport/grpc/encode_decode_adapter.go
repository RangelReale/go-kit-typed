package grpc

import (
	"context"

	"github.com/RangelReale/go-kit-typed/util"
	gokitgrpctransport "github.com/go-kit/kit/transport/grpc"
)

// DecodeRequestFuncReverseAdapter is an adapter tp the non-generic DecodeRequestFunc function
func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokitgrpctransport.DecodeRequestFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		return f(ctx, i)
	}
}

// EncodeRequestFuncReverseAdapter is an adapter to the non-generic EncodeRequestFunc function
func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokitgrpctransport.EncodeRequestFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		var req interface{}
		err := util.CallTypeWithError[Req](i, func(r Req) error {
			var callErr error
			req, callErr = f(ctx, r)
			return callErr
		})
		return req, err
	}
}

// EncodeResponseFuncReverseAdapter is an adapter to the non-generic EncodeResponseFunc function
func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitgrpctransport.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		var resp interface{}
		err := util.CallTypeWithError[Resp](i, func(r Resp) error {
			var callErr error
			resp, callErr = f(ctx, r)
			return callErr
		})
		return resp, err
	}
}

// DecodeResponseFuncReverseAdapter is an adapter to the non-generic DecodeResponseFunc function
func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokitgrpctransport.DecodeResponseFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		return f(ctx, i)
	}
}
