package nats

import (
	"context"

	"github.com/nats-io/nats.go"
)

// DecodeRequestFunc extracts a user-domain request object from a publisher
// request object. It's designed to be used in NATS subscribers, for subscriber-side
// endpoints. One straightforward DecodeRequestFunc could be something that
// JSON decodes from the request body to the concrete response type.
type DecodeRequestFunc[Req any] func(context.Context, *nats.Msg) (request Req, err error)

// EncodeRequestFunc encodes the passed request object into the NATS request
// object. It's designed to be used in NATS publishers, for publisher-side
// endpoints. One straightforward EncodeRequestFunc could something that JSON
// encodes the object directly to the request payload.
type EncodeRequestFunc[Req any] func(context.Context, *nats.Msg, Req) error

// EncodeResponseFunc encodes the passed response object to the subscriber reply.
// It's designed to be used in NATS subscribers, for subscriber-side
// endpoints. One straightforward EncodeResponseFunc could be something that
// JSON encodes the object directly to the response body.
type EncodeResponseFunc[Resp any] func(context.Context, string, *nats.Conn, Resp) error

// DecodeResponseFunc extracts a user-domain response object from an NATS
// response object. It's designed to be used in NATS publisher, for publisher-side
// endpoints. One straightforward DecodeResponseFunc could be something that
// JSON decodes from the response payload to the concrete response type.
type DecodeResponseFunc[Resp any] func(context.Context, *nats.Msg) (response Resp, err error)
