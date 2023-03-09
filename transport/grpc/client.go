package grpc

import (
	"github.com/RangelReale/go-kit-typed/endpoint"
	gokitgrpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// Client wraps a URL and provides a method that implements endpoint.Endpoint.
type Client[Req any, Resp any] struct {
	client *gokitgrpctransport.Client
}

// NewClient constructs a usable Client for a single remote endpoint.
// Pass an zero-value protobuf message of the RPC response type as
// the grpcReply argument.
func NewClient[Req any, Resp any](
	cc *grpc.ClientConn,
	serviceName string,
	method string,
	enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp],
	grpcReply interface{},
	options ...gokitgrpctransport.ClientOption,
) *Client[Req, Resp] {
	client := gokitgrpctransport.NewClient(
		cc,
		serviceName,
		method,
		EncodeRequestFuncReverseAdapter(enc),
		DecodeResponseFuncReverseAdapter(dec),
		grpcReply,
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// Endpoint returns a usable endpoint that will invoke the gRPC specified by the
// client.
func (c Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return endpoint.Adapter[Req, Resp](c.client.Endpoint())
}
