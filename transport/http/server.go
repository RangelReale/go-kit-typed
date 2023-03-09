package http

import (
	"net/http"

	"github.com/RangelReale/go-kit-typed/endpoint"
	gokithttptransport "github.com/go-kit/kit/transport/http"
)

// Server wraps an endpoint and implements http.Handler.
type Server[Req any, Resp any] struct {
	server *gokithttptransport.Server
}

// NewServer constructs a new server, which implements http.Handler and wraps
// the provided endpoint.
func NewServer[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	options ...gokithttptransport.ServerOption,
) *Server[Req, Resp] {
	server := gokithttptransport.NewServer(
		endpoint.ReverseAdapter(e),
		DecodeRequestFuncReverseAdapter(dec),
		EncodeResponseFuncReverseAdapter(enc),
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// NewServerStdDec constructs a new server, which implements http.Handler and wraps
// the provided endpoint, using the non-typed decoder.
func NewServerStdDec[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec gokithttptransport.DecodeRequestFunc,
	enc EncodeResponseFunc[Resp],
	options ...gokithttptransport.ServerOption,
) *Server[Req, Resp] {
	server := gokithttptransport.NewServer(
		endpoint.ReverseAdapter(e),
		dec,
		EncodeResponseFuncReverseAdapter(enc),
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// NewServerStdEnc constructs a new server, which implements http.Handler and wraps
// the provided endpoint, using the non-typed encoder.
func NewServerStdEnc[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc gokithttptransport.EncodeResponseFunc,
	options ...gokithttptransport.ServerOption,
) *Server[Req, Resp] {
	server := gokithttptransport.NewServer(
		endpoint.ReverseAdapter(e),
		DecodeRequestFuncReverseAdapter(dec),
		enc,
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// ServeHTTP implements http.Handler.
func (s Server[Req, Resp]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
