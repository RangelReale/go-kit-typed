package nats

import (
	"context"

	"github.com/RangelReale/go-kit-typed/util"
	gokitnatstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

// DecodeRequestFuncReverseAdapter is an adapter tp the non-generic DecodeRequestFunc function
func DecodeRequestFuncReverseAdapter[Req any](f DecodeRequestFunc[Req]) gokitnatstransport.DecodeRequestFunc {
	return func(ctx context.Context, msg *nats.Msg) (interface{}, error) {
		return f(ctx, msg)
	}
}

// EncodeRequestFuncReverseAdapter is an adapter to the non-generic EncodeRequestFunc function
func EncodeRequestFuncReverseAdapter[Req any](f EncodeRequestFunc[Req]) gokitnatstransport.EncodeRequestFunc {
	return func(ctx context.Context, msg *nats.Msg, i interface{}) error {
		return util.CallTypeWithError[Req](i, func(r Req) error {
			return f(ctx, msg, r)
		})
	}
}

// EncodeResponseFuncReverseAdapter is an adapter to the non-generic EncodeResponseFunc function
func EncodeResponseFuncReverseAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitnatstransport.EncodeResponseFunc {
	return func(ctx context.Context, s string, conn *nats.Conn, i interface{}) error {
		return util.CallTypeWithError[Resp](i, func(r Resp) error {
			return f(ctx, s, conn, r)
		})
	}
}

// DecodeResponseFuncReverseAdapter is an adapter to the non-generic DecodeResponseFunc function
func DecodeResponseFuncReverseAdapter[Resp any](f DecodeResponseFunc[Resp]) gokitnatstransport.DecodeResponseFunc {
	return func(ctx context.Context, msg *nats.Msg) (response interface{}, err error) {
		return f(ctx, msg)
	}
}
