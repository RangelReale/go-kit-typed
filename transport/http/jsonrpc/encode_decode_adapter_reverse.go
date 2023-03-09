package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

// EndpointCodecReverseAdapter is an adapter to the non-generic EndpointCodec type
func EndpointCodecReverseAdapter[Req any, Resp any](codec EndpointCodec[Req, Resp]) gokitjsonrpctransport.EndpointCodec {
	return gokitjsonrpctransport.EndpointCodec{
		Endpoint: endpoint.ReverseAdapter[Req, Resp](codec.Endpoint),
		Decode:   DecodeRequestFuncReverseAdapter(codec.Decode),
		Encode:   EncodeResponseFuncReverseAdapter(codec.Encode),
	}
}

// DecodeRequestFuncReverseAdapter is an adapter to the non-generic DecodeRequestFunc function
func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokitjsonrpctransport.DecodeRequestFunc {
	return func(ctx context.Context, message json.RawMessage) (interface{}, error) {
		return f(ctx, message)
	}
}

// EncodeResponseFuncReverseAdapter is an adapter to the non-generic EncodeResponseFunc function
func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitjsonrpctransport.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
		return util.CallTypeResponseWithError[Resp, json.RawMessage](i, func(r Resp) (json.RawMessage, error) {
			return f(ctx, r)
		})
	}
}

// EncodeRequestFuncReverseAdapter is an adapter to the non-generic EncodeRequestFunc function
func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokitjsonrpctransport.EncodeRequestFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
		return util.CallTypeResponseWithError[Req, json.RawMessage](i, func(r Req) (json.RawMessage, error) {
			return f(ctx, r)
		})
	}
}

// DecodeResponseFuncReverseAdapter is an adapter to the non-generic DecodeResponseFunc function
func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokitjsonrpctransport.DecodeResponseFunc {
	return func(ctx context.Context, response gokitjsonrpctransport.Response) (interface{}, error) {
		return f(ctx, response)
	}
}
