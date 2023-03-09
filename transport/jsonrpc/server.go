package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/RangelReale/go-kit-typed/endpoint"
	"github.com/RangelReale/go-kit-typed/util"
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

func serverEndpointCodecAdapter[Req any, Resp any](codec EndpointCodec[Req, Resp]) gokitjsonrpctransport.EndpointCodec {
	return gokitjsonrpctransport.EndpointCodec{
		Endpoint: endpoint.ReverseAdapter(codec.Endpoint),
		Decode:   serverDecodeRequestFuncAdapter(codec.Decode),
		Encode:   serverEncodeResponseFuncAdapter(codec.Encode),
	}
}

// func serverEndpointCodecMapAdapter[Req any, Resp any](ecm EndpointCodecMap[Req, Resp]) gokitjsonrpctransport.EndpointCodecMap {
// 	recm := make(gokitjsonrpctransport.EndpointCodecMap)
// 	for en, ev := range ecm {
// 		recm[en] = serverEndpointCodecAdapter(ev)
// 	}
// 	return recm
// }

func serverDecodeRequestFuncAdapter[Req any](f DecodeRequestFunc[Req]) gokitjsonrpctransport.DecodeRequestFunc {
	return func(ctx context.Context, message json.RawMessage) (interface{}, error) {
		return f(ctx, message)
	}
}

func serverEncodeResponseFuncAdapter[Resp any](f EncodeResponseFunc[Resp]) gokitjsonrpctransport.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (json.RawMessage, error) {
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
