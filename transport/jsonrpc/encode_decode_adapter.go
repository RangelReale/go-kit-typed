package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

// EndpointCodecAdapter is an adapter from the non-generic EndpointCodec type
func EndpointCodecAdapter[Req any, Resp any](codec gokitjsonrpctransport.EndpointCodec) EndpointCodec[Req, Resp] {
	return EndpointCodec[Req, Resp]{
		Endpoint: endpoint.Adapter[Req, Resp](codec.Endpoint),
		Decode:   DecodeRequestFuncAdapter[Req](codec.Decode),
		Encode:   EncodeResponseFuncAdapter[Resp](codec.Encode),
	}
}

// EndpointCodecReverseAdapter is an adapter to the non-generic EndpointCodec type
func EndpointCodecReverseAdapter[Req any, Resp any](codec EndpointCodec[Req, Resp]) gokitjsonrpctransport.EndpointCodec {
	return gokitjsonrpctransport.EndpointCodec{
		Endpoint: endpoint.ReverseAdapter[Req, Resp](codec.Endpoint),
		Decode:   DecodeRequestFuncReverseAdapter(codec.Decode),
		Encode:   EncodeResponseFuncReverseAdapter(codec.Encode),
	}
}

// DecodeRequestFuncAdapter is an adapter from the non-generic DecodeRequestFunc function
func DecodeRequestFuncAdapter[Req any](f gokitjsonrpctransport.DecodeRequestFunc) DecodeRequestFunc[Req] {
	return func(ctx context.Context, message json.RawMessage) (Req, error) {
		return util.ReturnTypeWithError[Req](f(ctx, message))
	}
}

// DecodeRequestFuncReverseAdapter is an adapter to the non-generic DecodeRequestFunc function
func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokitjsonrpctransport.DecodeRequestFunc {
	return func(ctx context.Context, message json.RawMessage) (interface{}, error) {
		return f(ctx, message)
	}
}

// EncodeResponseFuncAdapter is an adapter from the non-generic EncodeResponseFunc function
func EncodeResponseFuncAdapter[Resp any](f gokitjsonrpctransport.EncodeResponseFunc) EncodeResponseFunc[Resp] {
	return func(ctx context.Context, resp Resp) (json.RawMessage, error) {
		return f(ctx, resp)
	}
}

// EncodeResponseFuncReverseAdapter is an adapter to the non-generic EncodeResponseFunc function
func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitjsonrpctransport.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
		var rm json.RawMessage
		err := util.CallTypeWithError[Resp](i, func(r Resp) error {
			var callErr error
			rm, callErr = f(ctx, r)
			return callErr
		})
		return rm, err
	}
}

// EncodeRequestFuncReverseAdapter is an adapter to the non-generic EncodeRequestFunc function
func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokitjsonrpctransport.EncodeRequestFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
		var rm json.RawMessage
		err := util.CallTypeWithError[Req](i, func(r Req) error {
			var callErr error
			rm, callErr = f(ctx, r)
			return callErr
		})
		return rm, err
	}
}

// DecodeResponseFuncReverseAdapter is an adapter to the non-generic DecodeResponseFunc function
func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokitjsonrpctransport.DecodeResponseFunc {
	return func(ctx context.Context, response gokitjsonrpctransport.Response) (interface{}, error) {
		return f(ctx, response)
	}
}
