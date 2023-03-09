package endpoint

import (
	"context"
)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint[Req any, Resp any] func(ctx context.Context, request Req) (response Resp, err error)

// Middleware is a chainable behavior modifier for endpoints.
type Middleware[Req any, Resp any] func(Endpoint[Req, Resp]) Endpoint[Req, Resp]
