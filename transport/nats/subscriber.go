package nats

import (
	"github.com/RangelReale/go-kit-typed/endpoint"
	gokitnatstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

// Subscriber wraps an endpoint and provides nats.MsgHandler.
type Subscriber[Req any, Resp any] struct {
	subscriber *gokitnatstransport.Subscriber
}

// NewSubscriber constructs a new subscriber, which provides nats.MsgHandler and wraps
// the provided endpoint.
func NewSubscriber[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	options ...gokitnatstransport.SubscriberOption,
) *Subscriber[Req, Resp] {
	subscriber := gokitnatstransport.NewSubscriber(
		endpoint.ReverseAdapter(e),
		DecodeRequestFuncReverseAdapter(dec),
		EncodeResponseFuncReverseAdapter(enc),
		options...)
	return &Subscriber[Req, Resp]{
		subscriber: subscriber,
	}
}

// ServeMsg provides nats.MsgHandler.
func (s Subscriber[Req, Resp]) ServeMsg(nc *nats.Conn) func(msg *nats.Msg) {
	return s.subscriber.ServeMsg(nc)
}
