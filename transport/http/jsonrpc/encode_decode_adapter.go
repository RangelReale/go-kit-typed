package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

// MakeEndpointCodec creates a standard EndpointCodec from generic parameters.
// This is intended to be used when build manually a EndpointCodecMap.
func MakeEndpointCodec[Req any, Resp any](endpoint endpoint.Endpoint[Req, Resp], dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp]) gokitjsonrpctransport.EndpointCodec {
	return EndpointCodecReverseAdapter(EndpointCodec[Req, Resp]{
		Endpoint: endpoint,
		Decode:   dec,
		Encode:   enc,
	})
}

// EndpointCodecAdapter is an adapter from the non-generic EndpointCodec type
func EndpointCodecAdapter[Req any, Resp any](codec gokitjsonrpctransport.EndpointCodec) EndpointCodec[Req, Resp] {
	return EndpointCodec[Req, Resp]{
		Endpoint: endpoint.Adapter[Req, Resp](codec.Endpoint),
		Decode:   DecodeRequestFuncAdapter[Req](codec.Decode),
		Encode:   EncodeResponseFuncAdapter[Resp](codec.Encode),
	}
}

// DecodeRequestFuncAdapter is an adapter from the non-generic DecodeRequestFunc function
func DecodeRequestFuncAdapter[Req any](f gokitjsonrpctransport.DecodeRequestFunc) DecodeRequestFunc[Req] {
	return func(ctx context.Context, message json.RawMessage) (Req, error) {
		return util.ReturnTypeWithError[Req](f(ctx, message))
	}
}

// EncodeResponseFuncAdapter is an adapter from the non-generic EncodeResponseFunc function
func EncodeResponseFuncAdapter[Resp any](f gokitjsonrpctransport.EncodeResponseFunc) EncodeResponseFunc[Resp] {
	return func(ctx context.Context, resp Resp) (json.RawMessage, error) {
		return f(ctx, resp)
	}
}
