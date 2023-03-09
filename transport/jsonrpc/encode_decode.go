package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

// Server-Side Codec

// EndpointCodec defines a server Endpoint and its associated codecs
type EndpointCodec[Req any, Resp any] struct {
	Endpoint endpoint.Endpoint[Req, Resp]
	Decode   DecodeRequestFunc[Req]
	Encode   EncodeResponseFunc[Resp]
}

// // EndpointCodecMap maps the Request.Method to the proper EndpointCodec
// type EndpointCodecMap[Req any, Resp any] map[string]EndpointCodec[Req, Resp]

// DecodeRequestFunc extracts a user-domain request object from raw JSON
// It's designed to be used in JSON RPC servers, for server-side endpoints.
// One straightforward DecodeRequestFunc could be something that unmarshals
// JSON from the request body to the concrete request type.
type DecodeRequestFunc[Req any] func(context.Context, json.RawMessage) (request Req, err error)

// EncodeResponseFunc encodes the passed response object to a JSON RPC result.
// It's designed to be used in HTTP servers, for server-side endpoints.
// One straightforward EncodeResponseFunc could be something that JSON encodes
// the object directly.
type EncodeResponseFunc[Resp any] func(context.Context, Resp) (response json.RawMessage, err error)

// Client-Side Codec

// EncodeRequestFunc encodes the given request object to raw JSON.
// It's designed to be used in JSON RPC clients, for client-side
// endpoints. One straightforward EncodeResponseFunc could be something that
// JSON encodes the object directly.
type EncodeRequestFunc[Req any] func(context.Context, Req) (request json.RawMessage, err error)

// DecodeResponseFunc extracts a user-domain response object from an JSON RPC
// response object. It's designed to be used in JSON RPC clients, for
// client-side endpoints. It is the responsibility of this function to decide
// whether any error present in the JSON RPC response should be surfaced to the
// client endpoint.
type DecodeResponseFunc[Resp any] func(context.Context, gokitjsonrpctransport.Response) (response Resp, err error)

func EndpointCodecAdapter[Req any, Resp any](codec gokitjsonrpctransport.EndpointCodec) EndpointCodec[Req, Resp] {
	return EndpointCodec[Req, Resp]{
		Endpoint: endpoint.Adapter[Req, Resp](codec.Endpoint),
		Decode:   DecodeRequestFuncAdapter[Req](codec.Decode),
		Encode:   EncodeResponseFuncAdapter[Resp](codec.Encode),
	}
}

func EndpointCodecReverseAdapter[Req any, Resp any](codec EndpointCodec[Req, Resp]) gokitjsonrpctransport.EndpointCodec {
	return gokitjsonrpctransport.EndpointCodec{
		Endpoint: endpoint.ReverseAdapter[Req, Resp](codec.Endpoint),
		Decode:   DecodeRequestFuncReverseAdapter(codec.Decode),
		Encode:   EncodeResponseFuncReverseAdapter(codec.Encode),
	}
}

func DecodeRequestFuncAdapter[Req any](f gokitjsonrpctransport.DecodeRequestFunc) DecodeRequestFunc[Req] {
	return func(ctx context.Context, message json.RawMessage) (Req, error) {
		req, err := f(ctx, message)
		if err != nil {
			var rr Req
			return rr, err
		}

		switch tr := req.(type) {
		case nil:
			var rr Req
			return rr, nil
		case Req:
			return tr, nil
		default:
			var rr Req
			return rr, util.ErrParameterInvalidType
		}
	}
}

func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokitjsonrpctransport.DecodeRequestFunc {
	return func(ctx context.Context, message json.RawMessage) (interface{}, error) {
		return f(ctx, message)
	}
}

func EncodeResponseFuncAdapter[Resp any](f gokitjsonrpctransport.EncodeResponseFunc) EncodeResponseFunc[Resp] {
	return func(ctx context.Context, resp Resp) (json.RawMessage, error) {
		return f(ctx, resp)
	}
}

func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitjsonrpctransport.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
		switch ti := i.(type) {
		case nil:
			var r Resp
			return f(ctx, r)
		case Resp:
			return f(ctx, ti)
		default:
			return nil, util.ErrParameterInvalidType
		}
	}
}

func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokitjsonrpctransport.EncodeRequestFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
		switch ri := i.(type) {
		case nil:
			var r Req
			return f(ctx, r)
		case Req:
			return f(ctx, ri)
		default:
			return nil, util.ErrParameterInvalidType
		}
	}
}

func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokitjsonrpctransport.DecodeResponseFunc {
	return func(ctx context.Context, response gokitjsonrpctransport.Response) (interface{}, error) {
		return f(ctx, response)
	}
}
