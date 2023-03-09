package nats_test

import (
	"context"
	"strings"
	"testing"
	"time"

	natstransport "github.com/RangelReale/go-kit-typed/transport/nats"
	gokitnatstransport "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

type clientReq struct {
	req string
}

// type clientResp struct {
// 	resp string
// }

func TestPublisher(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, clientReq) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (TestResponse, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	s, c := newNATSConn(t)
	defer func() { s.Shutdown(); s.WaitForShutdown() }()
	defer c.Close()

	sub, err := c.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := c.Publish(msg.Reply, []byte(testdata)); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		c,
		"natstransport.test",
		encode,
		decode,
	)

	response, err := publisher.Endpoint()(context.Background(), clientReq{"req1"})
	if err != nil {
		t.Fatal(err)
	}

	if want, have := testdata, response.String; want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}

func TestPublisherBefore(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, clientReq) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (TestResponse, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	s, c := newNATSConn(t)
	defer func() { s.Shutdown(); s.WaitForShutdown() }()
	defer c.Close()

	sub, err := c.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := c.Publish(msg.Reply, msg.Data); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		c,
		"natstransport.test",
		encode,
		decode,
		gokitnatstransport.PublisherBefore(func(ctx context.Context, msg *nats.Msg) context.Context {
			msg.Data = []byte(strings.ToUpper(string(testdata)))
			return ctx
		}),
	)

	response, err := publisher.Endpoint()(context.Background(), clientReq{"req1"})
	if err != nil {
		t.Fatal(err)
	}

	if want, have := strings.ToUpper(testdata), response.String; want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}

func TestPublisherAfter(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, clientReq) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (TestResponse, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	s, c := newNATSConn(t)
	defer func() { s.Shutdown(); s.WaitForShutdown() }()
	defer c.Close()

	sub, err := c.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := c.Publish(msg.Reply, []byte(testdata)); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		c,
		"natstransport.test",
		encode,
		decode,
		gokitnatstransport.PublisherAfter(func(ctx context.Context, msg *nats.Msg) context.Context {
			msg.Data = []byte(strings.ToUpper(string(msg.Data)))
			return ctx
		}),
	)

	response, err := publisher.Endpoint()(context.Background(), clientReq{"req1"})
	if err != nil {
		t.Fatal(err)
	}

	if want, have := strings.ToUpper(testdata), response.String; want != have {
		t.Errorf("want %q, have %q", want, have)
	}

}

func TestPublisherTimeout(t *testing.T) {
	var (
		encode = func(context.Context, *nats.Msg, clientReq) error { return nil }
		decode = func(_ context.Context, msg *nats.Msg) (TestResponse, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	s, c := newNATSConn(t)
	defer func() { s.Shutdown(); s.WaitForShutdown() }()
	defer c.Close()

	ch := make(chan struct{})
	defer close(ch)

	sub, err := c.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		<-ch
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		c,
		"natstransport.test",
		encode,
		decode,
		gokitnatstransport.PublisherTimeout(time.Second),
	)

	_, err = publisher.Endpoint()(context.Background(), clientReq{"req1"})
	if err != context.DeadlineExceeded {
		t.Errorf("want %s, have %s", context.DeadlineExceeded, err)
	}
}

func TestPublisherCancellation(t *testing.T) {
	var (
		testdata = "testdata"
		encode   = func(context.Context, *nats.Msg, clientReq) error { return nil }
		decode   = func(_ context.Context, msg *nats.Msg) (TestResponse, error) {
			return TestResponse{string(msg.Data), ""}, nil
		}
	)

	s, c := newNATSConn(t)
	defer func() { s.Shutdown(); s.WaitForShutdown() }()
	defer c.Close()

	sub, err := c.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		if err := c.Publish(msg.Reply, []byte(testdata)); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		c,
		"natstransport.test",
		encode,
		decode,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = publisher.Endpoint()(ctx, clientReq{"req1"})
	if err != context.Canceled {
		t.Errorf("want %s, have %s", context.Canceled, err)
	}
}

func TestEncodeJSONRequest(t *testing.T) {
	var data string

	s, c := newNATSConn(t)
	defer func() { s.Shutdown(); s.WaitForShutdown() }()
	defer c.Close()

	sub, err := c.QueueSubscribe("natstransport.test", "natstransport", func(msg *nats.Msg) {
		data = string(msg.Data)

		if err := c.Publish(msg.Reply, []byte("")); err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	publisher := natstransport.NewPublisher(
		c,
		"natstransport.test",
		gokitnatstransport.EncodeJSONRequest,
		func(context.Context, *nats.Msg) (interface{}, error) { return nil, nil },
	).Endpoint()

	for _, test := range []struct {
		value interface{}
		body  string
	}{
		{nil, "null"},
		{12, "12"},
		{1.2, "1.2"},
		{true, "true"},
		{"test", "\"test\""},
		{struct {
			Foo string `json:"foo"`
		}{"foo"}, "{\"foo\":\"foo\"}"},
	} {
		if _, err := publisher(context.Background(), test.value); err != nil {
			t.Fatal(err)
			continue
		}

		if data != test.body {
			t.Errorf("%v: actual %#v, expected %#v", test.value, data, test.body)
		}
	}

}
