package http

import (
	"net/url"

	"github.com/RangelReale/go-kit-typed/endpoint"
	gokithttptransport "github.com/go-kit/kit/transport/http"
)

// Client wraps a URL and provides a method that implements endpoint.Endpoint.
type Client[Req any, Resp any] struct {
	client *gokithttptransport.Client
}

// NewClient constructs a usable Client for a single remote method.
func NewClient[Req any, Resp any](method string, tgt *url.URL, enc EncodeRequestFunc[Req],
	dec DecodeResponseFunc[Resp], options ...gokithttptransport.ClientOption) *Client[Req, Resp] {
	client := gokithttptransport.NewClient(
		method,
		tgt,
		EncodeRequestFuncReverseAdapter(enc),
		DecodeResponseFuncReverseAdapter(dec),
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// NewClientStdEnc constructs a usable Client for a single remote method.
func NewClientStdEnc[Req any, Resp any](method string, tgt *url.URL, enc gokithttptransport.EncodeRequestFunc,
	dec DecodeResponseFunc[Resp], options ...gokithttptransport.ClientOption) *Client[Req, Resp] {
	client := gokithttptransport.NewClient(method,
		tgt,
		enc,
		DecodeResponseFuncReverseAdapter(dec),
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// NewClientStdDec constructs a usable Client for a single remote method.
func NewClientStdDec[Req any, Resp any](method string, tgt *url.URL, enc EncodeRequestFunc[Req],
	dec gokithttptransport.DecodeResponseFunc, options ...gokithttptransport.ClientOption) *Client[Req, Resp] {
	client := gokithttptransport.NewClient(
		method,
		tgt,
		EncodeRequestFuncReverseAdapter(enc),
		dec,
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// NewExplicitClient is like NewClient but uses a CreateRequestFunc instead of a
// method, target URL, and EncodeRequestFunc, which allows for more control over
// the outgoing HTTP request.
func NewExplicitClient[Req any, Resp any](req CreateRequestFunc[Req], dec DecodeResponseFunc[Resp],
	options ...gokithttptransport.ClientOption) *Client[Req, Resp] {
	client := gokithttptransport.NewExplicitClient(
		CreateRequestFuncReverseAdapter(req),
		DecodeResponseFuncReverseAdapter(dec),
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// NewExplicitClientStdCreate is like NewClient but uses a CreateRequestFunc instead of a
// method, target URL, and EncodeRequestFunc, which allows for more control over
// the outgoing HTTP request, using the non-typed creator.
func NewExplicitClientStdCreate[Req any, Resp any](req gokithttptransport.CreateRequestFunc, dec DecodeResponseFunc[Resp],
	options ...gokithttptransport.ClientOption) *Client[Req, Resp] {
	client := gokithttptransport.NewExplicitClient(
		req,
		DecodeResponseFuncReverseAdapter(dec),
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// NewExplicitClientStdDec is like NewClient but uses a CreateRequestFunc instead of a
// method, target URL, and EncodeRequestFunc, which allows for more control over
// the outgoing HTTP request, using the non-typed decoder.
func NewExplicitClientStdDec[Req any, Resp any](req CreateRequestFunc[Req], dec gokithttptransport.DecodeResponseFunc,
	options ...gokithttptransport.ClientOption) *Client[Req, Resp] {
	client := gokithttptransport.NewExplicitClient(CreateRequestFuncReverseAdapter(req),
		dec,
		options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// Endpoint returns a usable Go kit endpoint that calls the remote HTTP endpoint.
func (c Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return endpoint.Adapter[Req, Resp](c.client.Endpoint())
}
