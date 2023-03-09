package jsonrpc

import (
	"net/url"

	"github.com/RangelReale/go-kit-typed/endpoint"
	gokithttptransport "github.com/go-kit/kit/transport/http"
	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

// Client wraps a URL and provides a method that implements endpoint.Endpoint.
type Client[Req any, Resp any] struct {
	client *gokitjsonrpctransport.Client
}

// NewClient constructs a usable Client for a single remote method.
func NewClient[Req any, Resp any](
	tgt *url.URL,
	method string,
	options ...ClientOption,
) *Client[Req, Resp] {
	var copt clientOptions
	for _, opt := range options {
		opt(&copt)
	}
	client := gokitjsonrpctransport.NewClient(tgt, method, copt.options...)
	return &Client[Req, Resp]{
		client: client,
	}
}

// Endpoint returns a usable Go kit endpoint that calls the remote HTTP endpoint.
func (c Client[Req, Resp]) Endpoint() endpoint.Endpoint[Req, Resp] {
	return endpoint.Adapter[Req, Resp](c.client.Endpoint())
}

type clientOptions struct {
	options []gokitjsonrpctransport.ClientOption
}

// ClientOption sets an optional parameter for clients.
type ClientOption func(*clientOptions)

// SetClient sets the underlying HTTP client used for requests.
// By default, http.DefaultClient is used.
func SetClient(client gokithttptransport.HTTPClient) ClientOption {
	return func(c *clientOptions) { c.options = append(c.options, gokitjsonrpctransport.SetClient(client)) }
}

// ClientBefore sets the RequestFuncs that are applied to the outgoing HTTP
// request before it's invoked.
func ClientBefore(before ...gokithttptransport.RequestFunc) ClientOption {
	return func(c *clientOptions) { c.options = append(c.options, gokitjsonrpctransport.ClientBefore(before...)) }
}

// ClientAfter sets the ClientResponseFuncs applied to the server's HTTP
// response prior to it being decoded. This is useful for obtaining anything
// from the response and adding onto the context prior to decoding.
func ClientAfter(after ...gokithttptransport.ClientResponseFunc) ClientOption {
	return func(c *clientOptions) { c.options = append(c.options, gokitjsonrpctransport.ClientAfter(after...)) }
}

// ClientFinalizer is executed at the end of every HTTP request.
// By default, no finalizer is registered.
func ClientFinalizer(f gokithttptransport.ClientFinalizerFunc) ClientOption {
	return func(c *clientOptions) { c.options = append(c.options, gokitjsonrpctransport.ClientFinalizer(f)) }
}

// ClientRequestEncoder sets the func used to encode the request params to JSON.
// If not set, DefaultRequestEncoder is used.
func ClientRequestEncoder[Req any](enc EncodeRequestFunc[Req]) ClientOption {
	return func(c *clientOptions) {
		c.options = append(c.options,
			gokitjsonrpctransport.ClientRequestEncoder(EncodeRequestFuncReverseAdapter(enc)))
	}
}

// ClientResponseDecoder sets the func used to decode the response params from
// JSON. If not set, DefaultResponseDecoder is used.
func ClientResponseDecoder[Req any](dec DecodeResponseFunc[Req]) ClientOption {
	return func(c *clientOptions) {
		c.options = append(c.options,
			gokitjsonrpctransport.ClientResponseDecoder(DecodeResponseFuncReverseAdapter(dec)))
	}
}

// ClientRequestIDGenerator is executed before each request to generate an ID
// for the request.
// By default, AutoIncrementRequestID is used.
func ClientRequestIDGenerator(g gokitjsonrpctransport.RequestIDGenerator) ClientOption {
	return func(c *clientOptions) {
		c.options = append(c.options, gokitjsonrpctransport.ClientRequestIDGenerator(g))
	}
}

// BufferedStream sets whether the Response.Body is left open, allowing it
// to be read from later. Useful for transporting a file as a buffered stream.
func BufferedStream(buffered bool) ClientOption {
	return func(c *clientOptions) { c.options = append(c.options, gokitjsonrpctransport.BufferedStream(buffered)) }
}
