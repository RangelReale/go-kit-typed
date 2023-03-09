package nats

import (
	"github.com/RangelReale/go-kit-typed/endpoint"
	gokitnatstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

// Publisher wraps a URL and provides a method that implements endpoint.Endpoint.
type Publisher[Req any, Resp any] struct {
	publisher *gokitnatstransport.Publisher
}

// NewPublisher constructs a usable Publisher for a single remote method.
func NewPublisher[Req any, Resp any](
	publisher *nats.Conn,
	subject string,
	enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp],
	options ...gokitnatstransport.PublisherOption,
) *Publisher[Req, Resp] {
	pb := gokitnatstransport.NewPublisher(
		publisher,
		subject,
		EncodeRequestFuncReverseAdapter(enc),
		DecodeResponseFuncReverseAdapter(dec),
		options...)
	return &Publisher[Req, Resp]{
		publisher: pb,
	}
}

// NewPublisherStdEnc constructs a usable Publisher for a single remote method.
func NewPublisherStdEnc[Req any, Resp any](
	publisher *nats.Conn,
	subject string,
	enc gokitnatstransport.EncodeRequestFunc,
	dec DecodeResponseFunc[Resp],
	options ...gokitnatstransport.PublisherOption,
) *Publisher[Req, Resp] {
	pb := gokitnatstransport.NewPublisher(
		publisher,
		subject,
		enc,
		DecodeResponseFuncReverseAdapter(dec),
		options...)
	return &Publisher[Req, Resp]{
		publisher: pb,
	}
}

// NewPublisherStdDec constructs a usable Publisher for a single remote method.
func NewPublisherStdDec[Req any, Resp any](
	publisher *nats.Conn,
	subject string,
	enc EncodeRequestFunc[Req],
	dec gokitnatstransport.DecodeResponseFunc,
	options ...gokitnatstransport.PublisherOption,
) *Publisher[Req, Resp] {
	pb := gokitnatstransport.NewPublisher(
		publisher,
		subject,
		EncodeRequestFuncReverseAdapter(enc),
		dec,
		options...)
	return &Publisher[Req, Resp]{
		publisher: pb,
	}
}

// Endpoint returns a usable endpoint that invokes the remote endpoint.
func (p Publisher[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return endpoint.Adapter[Req, Resp](p.publisher.Endpoint())
}
