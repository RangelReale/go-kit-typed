package grpc

import (
	"context"

	"github.com/RangelReale/go-kit-typed/endpoint"
	gokitgrpctransport "github.com/go-kit/kit/transport/grpc"
)

// Server wraps an endpoint and implements grpc.Handler.
type Server[Req any, Resp any] struct {
	server *gokitgrpctransport.Server
}

// NewServer constructs a new server, which implements wraps the provided
// endpoint and implements the Handler interface. Consumers should write
// bindings that adapt the concrete gRPC methods from their compiled protobuf
// definitions to individual handlers. Request and response objects are from the
// caller business domain, not gRPC request and reply types.
func NewServer[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	options ...gokitgrpctransport.ServerOption,
) *Server[Req, Resp] {
	server := gokitgrpctransport.NewServer(
		endpoint.ReverseAdapter(e),
		DecodeRequestFuncReverseAdapter(dec),
		EncodeResponseFuncReverseAdapter(enc),
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// NewServerStdDec constructs a new server, which implements wraps the provided
// endpoint and implements the Handler interface. Consumers should write
// bindings that adapt the concrete gRPC methods from their compiled protobuf
// definitions to individual handlers. Request and response objects are from the
// caller business domain, not gRPC request and reply types, using the non-typed decoder.
func NewServerStdDec[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec gokitgrpctransport.DecodeRequestFunc,
	enc EncodeResponseFunc[Resp],
	options ...gokitgrpctransport.ServerOption,
) *Server[Req, Resp] {
	server := gokitgrpctransport.NewServer(
		endpoint.ReverseAdapter(e),
		dec,
		EncodeResponseFuncReverseAdapter(enc),
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// NewServerStdEnc constructs a new server, which implements wraps the provided
// endpoint and implements the Handler interface. Consumers should write
// bindings that adapt the concrete gRPC methods from their compiled protobuf
// definitions to individual handlers. Request and response objects are from the
// caller business domain, not gRPC request and reply types, using the non-typed encoder.
func NewServerStdEnc[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc gokitgrpctransport.EncodeResponseFunc,
	options ...gokitgrpctransport.ServerOption,
) *Server[Req, Resp] {
	server := gokitgrpctransport.NewServer(
		endpoint.ReverseAdapter(e),
		DecodeRequestFuncReverseAdapter(dec),
		enc,
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// ServeGRPC implements the Handler interface.
func (s Server[Req, Resp]) ServeGRPC(ctx context.Context, req interface{}) (retctx context.Context, resp interface{}, err error) {
	return s.ServeGRPC(ctx, req)
}
