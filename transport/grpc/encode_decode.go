package grpc

import (
	"context"

	"github.com/RangelReale/go-kit-typed/util"
	gokitgrpctransport "github.com/go-kit/kit/transport/grpc"
)

// DecodeRequestFunc extracts a user-domain request object from a gRPC request.
// It's designed to be used in gRPC servers, for server-side endpoints. One
// straightforward DecodeRequestFunc could be something that decodes from the
// gRPC request message to the concrete request type.
type DecodeRequestFunc[Req any] func(context.Context, interface{}) (request Req, err error)

// EncodeRequestFunc encodes the passed request object into the gRPC request
// object. It's designed to be used in gRPC clients, for client-side endpoints.
// One straightforward EncodeRequestFunc could something that encodes the object
// directly to the gRPC request message.
type EncodeRequestFunc[Req any] func(context.Context, Req) (request interface{}, err error)

// EncodeResponseFunc encodes the passed response object to the gRPC response
// message. It's designed to be used in gRPC servers, for server-side endpoints.
// One straightforward EncodeResponseFunc could be something that encodes the
// object directly to the gRPC response message.
type EncodeResponseFunc[Resp any] func(context.Context, Resp) (response interface{}, err error)

// DecodeResponseFunc extracts a user-domain response object from a gRPC
// response object. It's designed to be used in gRPC clients, for client-side
// endpoints. One straightforward DecodeResponseFunc could be something that
// decodes from the gRPC response message to the concrete response type.
type DecodeResponseFunc[Resp any] func(context.Context, interface{}) (response Resp, err error)

func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokitgrpctransport.DecodeRequestFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		return f(ctx, i)
	}
}

func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokitgrpctransport.EncodeRequestFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
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

func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitgrpctransport.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
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

func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokitgrpctransport.DecodeResponseFunc {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		return f(ctx, i)
	}
}
