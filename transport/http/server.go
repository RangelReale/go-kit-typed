package http

import (
	"context"
	"net/http"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
	gokitendpoint "github.com/go-kit/kit/endpoint"
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
	server := gokithttptransport.NewServer(serverEndpointAdapter(e),
		serverDecodeRequestFuncAdapter(dec),
		serverEncodeResponseFuncAdapter(enc),
		options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// ServeHTTP implements http.Handler.
func (s Server[Req, Resp]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}

func serverEndpointAdapter[Req any, Resp any](e endpoint.Endpoint[Req, Resp]) gokitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		switch tr := request.(type) {
		case nil:
			var r Req
			return e(ctx, r)
		case Req:
			return e(ctx, tr)
		default:
			var r Req
			return r, util.ErrParameterInvalidType
		}
	}
}

func serverDecodeRequestFuncAdapter[Req any](f DecodeRequestFunc[Req]) gokithttptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return f(ctx, r)
	}
}

func serverEncodeResponseFuncAdapter[Resp any](f EncodeResponseFunc[Resp]) gokithttptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, i interface{}) error {
		switch ti := i.(type) {
		case nil:
			var r Resp
			return f(ctx, w, r)
		case Resp:
			return f(ctx, w, ti)
		default:
			return util.ErrParameterInvalidType
		}
	}
}
