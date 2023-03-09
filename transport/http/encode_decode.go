package http

import (
	"context"
	"net/http"

	"github.com/RangelReale/go-kit-typed/util"
	gokithttptransport "github.com/go-kit/kit/transport/http"
)

// DecodeRequestFunc extracts a user-domain request object from an HTTP
// request object. It's designed to be used in HTTP servers, for server-side
// endpoints. One straightforward DecodeRequestFunc could be something that
// JSON decodes from the request body to the concrete request type.
type DecodeRequestFunc[Req any] func(context.Context, *http.Request) (request Req, err error)

// EncodeRequestFunc encodes the passed request object into the HTTP request
// object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward EncodeRequestFunc could be something that JSON
// encodes the object directly to the request body.
type EncodeRequestFunc[Req any] func(context.Context, *http.Request, Req) error

// CreateRequestFunc creates an outgoing HTTP request based on the passed
// request object. It's designed to be used in HTTP clients, for client-side
// endpoints. It's a more powerful version of EncodeRequestFunc, and can be used
// if more fine-grained control of the HTTP request is required.
type CreateRequestFunc[Req any] func(context.Context, Req) (*http.Request, error)

// EncodeResponseFunc encodes the passed response object to the HTTP response
// writer. It's designed to be used in HTTP servers, for server-side
// endpoints. One straightforward EncodeResponseFunc could be something that
// JSON encodes the object directly to the response body.
type EncodeResponseFunc[Resp any] func(context.Context, http.ResponseWriter, Resp) error

// DecodeResponseFunc extracts a user-domain response object from an HTTP
// response object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward DecodeResponseFunc could be something that
// JSON decodes from the response body to the concrete response type.
type DecodeResponseFunc[Resp any] func(context.Context, *http.Response) (response Resp, err error)

// DecodeRequestFuncAdapter is an adapter from the non-generic DecodeRequestFunc function
func DecodeRequestFuncAdapter[Req any](f gokithttptransport.DecodeRequestFunc) DecodeRequestFunc[Req] {
	return func(ctx context.Context, r *http.Request) (Req, error) {
		req, err := f(ctx, r)
		if err != nil {
			var rr Req
			return rr, err
		}

		switch tr := req.(type) {
		case nil:
			var rr Req
			return rr, nil
		case Req:
			return tr, nil
		default:
			var rr Req
			return rr, util.ErrParameterInvalidType
		}
	}
}

// EncodeRequestFuncAdapter is an adapter from the non-generic EncodeRequestFunc function
func EncodeRequestFuncAdapter[Req any](f gokithttptransport.EncodeRequestFunc) EncodeRequestFunc[Req] {
	return func(ctx context.Context, request *http.Request, req Req) error {
		return f(ctx, request, req)
	}
}

// EncodeResponseFuncAdapter is an adapter from the non-generic EncodeResponseFunc function
func EncodeResponseFuncAdapter[Resp any](f gokithttptransport.EncodeResponseFunc) EncodeResponseFunc[Resp] {
	return func(ctx context.Context, writer http.ResponseWriter, resp Resp) error {
		return f(ctx, writer, resp)
	}
}

// DecodeResponseFuncAdapter is an adapter from the non-generic DecodeRequestFunc function
func DecodeResponseFuncAdapter[Resp any](f gokithttptransport.DecodeResponseFunc) DecodeResponseFunc[Resp] {
	return func(ctx context.Context, r *http.Response) (Resp, error) {
		resp, err := f(ctx, r)
		if err != nil {
			var rr Resp
			return rr, err
		}

		switch tr := resp.(type) {
		case nil:
			var rr Resp
			return rr, nil
		case Resp:
			return tr, nil
		default:
			var rr Resp
			return rr, util.ErrParameterInvalidType
		}
	}
}
