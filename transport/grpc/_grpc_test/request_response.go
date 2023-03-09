package test

import (
	"context"

	"github.com/RangelReale/go-kit-typed/transport/grpc/_grpc_test/pb"
)

func encodeRequest(ctx context.Context, r TestRequest) (interface{}, error) {
	return &pb.TestRequest{A: r.A, B: r.B}, nil
}

func decodeRequest(ctx context.Context, req interface{}) (TestRequest, error) {
	r := req.(*pb.TestRequest)
	return TestRequest{A: r.A, B: r.B}, nil
}

func encodeResponse(ctx context.Context, r *TestResponse) (interface{}, error) {
	return &pb.TestResponse{V: r.V}, nil
}

func decodeResponse(ctx context.Context, resp interface{}) (*TestResponse, error) {
	r := resp.(*pb.TestResponse)
	return &TestResponse{V: r.V, Ctx: ctx}, nil
}
