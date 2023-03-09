package test

import (
	"context"

	"google.golang.org/grpc"

	"github.com/RangelReale/go-kit-typed/endpoint"
	grpctransport "github.com/RangelReale/go-kit-typed/transport/grpc"
	"github.com/RangelReale/go-kit-typed/transport/grpc/_grpc_test/pb"
	gokitgrpctransport "github.com/go-kit/kit/transport/grpc"
)

type clientBinding struct {
	test endpoint.Endpoint[TestRequest, *TestResponse]
}

func (c *clientBinding) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	r, err := c.test(ctx, TestRequest{A: a, B: b})
	if err != nil {
		return nil, "", err
	}
	return r.Ctx, r.V, nil
}

func NewClient(cc *grpc.ClientConn) Service {
	return &clientBinding{
		test: grpctransport.NewClient(
			cc,
			"pb.Test",
			"Test",
			encodeRequest,
			decodeResponse,
			&pb.TestResponse{},
			gokitgrpctransport.ClientBefore(
				injectCorrelationID,
			),
			gokitgrpctransport.ClientBefore(
				displayClientRequestHeaders,
			),
			gokitgrpctransport.ClientAfter(
				displayClientResponseHeaders,
				displayClientResponseTrailers,
			),
			gokitgrpctransport.ClientAfter(
				extractConsumedCorrelationID,
			),
		).Endpoint(),
	}
}
