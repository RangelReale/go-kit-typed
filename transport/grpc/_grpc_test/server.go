package test

import (
	"context"
	"fmt"

	"github.com/RangelReale/go-kit-typed/endpoint"
	grpctransport "github.com/RangelReale/go-kit-typed/transport/grpc"
	"github.com/RangelReale/go-kit-typed/transport/grpc/_grpc_test/pb"
	gokitgrpctransport "github.com/go-kit/kit/transport/grpc"
)

type service struct{}

func (service) Test(ctx context.Context, a string, b int64) (context.Context, string, error) {
	return nil, fmt.Sprintf("%s = %d", a, b), nil
}

func NewService() Service {
	return service{}
}

func makeTestEndpoint(svc Service) endpoint.Endpoint[TestRequest, *TestResponse] {
	return func(ctx context.Context, req TestRequest) (*TestResponse, error) {
		newCtx, v, err := svc.Test(ctx, req.A, req.B)
		return &TestResponse{
			V:   v,
			Ctx: newCtx,
		}, err
	}
}

type serverBinding struct {
	pb.UnimplementedTestServer

	test gokitgrpctransport.Handler
}

func (b *serverBinding) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	_, response, err := b.test.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.(*pb.TestResponse), nil
}

func NewBinding(svc Service) *serverBinding {
	return &serverBinding{
		test: grpctransport.NewServer(
			makeTestEndpoint(svc),
			decodeRequest,
			encodeResponse,
			gokitgrpctransport.ServerBefore(
				extractCorrelationID,
			),
			gokitgrpctransport.ServerBefore(
				displayServerRequestHeaders,
			),
			gokitgrpctransport.ServerAfter(
				injectResponseHeader,
				injectResponseTrailer,
				injectConsumedCorrelationID,
			),
			gokitgrpctransport.ServerAfter(
				displayServerResponseHeaders,
				displayServerResponseTrailers,
			),
		),
	}
}
