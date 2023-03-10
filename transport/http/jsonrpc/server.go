package jsonrpc

import (
	"net/http"

	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

// Server wraps an endpoint and implements http.Handler.
type Server[Req any, Resp any] struct {
	server *gokitjsonrpctransport.Server
}

// NewServer constructs a new server, which implements http.Server.
func NewServer[Req any, Resp any](
	ecm gokitjsonrpctransport.EndpointCodecMap,
	options ...gokitjsonrpctransport.ServerOption,
) *Server[Req, Resp] {
	server := gokitjsonrpctransport.NewServer(ecm, options...)
	return &Server[Req, Resp]{
		server: server,
	}
}

// ServeHTTP implements http.Handler.
func (s Server[Req, Resp]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
