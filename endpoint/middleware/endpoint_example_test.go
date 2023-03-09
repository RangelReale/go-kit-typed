package middleware_test

import (
	"context"
	"fmt"

	"github.com/RangelReale/go-kit-typed/endpoint/middleware"
	gokitendpoint "github.com/go-kit/kit/endpoint"
)

func ExampleChainGeneric() {
	m := gokitendpoint.Chain(
		annotate("first"),
		annotate("second"),
		annotate("third"),
	)

	e := middleware.Wrapper(m, myGenericEndpoint)

	if _, err := e(ctx, "data1"); err != nil {
		panic(err)
	}

	// Output:
	// first pre
	// second pre
	// third pre
	// my hendpoint data1!
	// third post
	// second post
	// first post
}

var (
	ctx = context.Background()
)

func annotate(s string) gokitendpoint.Middleware {
	return func(next gokitendpoint.Endpoint) gokitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			fmt.Println(s, "pre")
			defer fmt.Println(s, "post")
			return next(ctx, request)
		}
	}
}

func myGenericEndpoint(ctx context.Context, x string) (int, error) {
	fmt.Printf("my hendpoint %s!\n", x)
	return 12, nil
}
