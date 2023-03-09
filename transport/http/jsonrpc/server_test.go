package jsonrpc_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RangelReale/go-kit-typed/transport/http/jsonrpc"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	gokitjsonrpctransport "github.com/go-kit/kit/transport/http/jsonrpc"
)

func addBody() io.Reader {
	return body(`{"jsonrpc": "2.0", "method": "add", "params": [3, 2], "id": 1}`)
}

func body(in string) io.Reader {
	return strings.NewReader(in)
}

func unmarshalResponse(body []byte) (resp gokitjsonrpctransport.Response, err error) {
	err = json.Unmarshal(body, &resp)
	return
}

func expectErrorCode(t *testing.T, want int, body []byte) {
	t.Helper()

	r, err := unmarshalResponse(body)
	if err != nil {
		t.Fatalf("Can't decode response: %v (%s)", err, body)
	}
	if r.Error == nil {
		t.Fatalf("Expected error on response. Got none: %s", body)
	}
	if have := r.Error.Code; want != have {
		t.Fatalf("Unexpected error code. Want %d, have %d: %s", want, have, body)
	}
}

func expectValidRequestID(t *testing.T, want int, body []byte) {
	t.Helper()

	r, err := unmarshalResponse(body)
	if err != nil {
		t.Fatalf("Can't decode response: %v (%s)", err, body)
	}
	have, err := r.ID.Int()
	if err != nil {
		t.Fatalf("Can't get requestID in response. err=%s, body=%s", err, body)
	}
	if want != have {
		t.Fatalf("Request ID: want %d, have %d (%s)", want, have, body)
	}
}

func expectNilRequestID(t *testing.T, body []byte) {
	t.Helper()

	r, err := unmarshalResponse(body)
	if err != nil {
		t.Fatalf("Can't decode response: %v (%s)", err, body)
	}
	if r.ID != nil {
		t.Fatalf("Request ID: want nil, have %v", r.ID)
	}
}

func nopDecoder(context.Context, json.RawMessage) (interface{}, error) { return struct{}{}, nil }
func nopEncoder(context.Context, interface{}) (json.RawMessage, error) { return []byte("[]"), nil }

type mockLogger struct {
	Called   bool
	LastArgs []interface{}
}

func (l *mockLogger) Log(keyvals ...interface{}) error {
	l.Called = true
	l.LastArgs = append(l.LastArgs, keyvals)
	return nil
}

func TestServerBadDecode(t *testing.T) {
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: gokitendpoint.Nop,
			Decode:   func(context.Context, json.RawMessage) (interface{}, error) { return struct{}{}, errors.New("oof") },
			Encode:   nopEncoder,
		},
	}
	logger := mockLogger{}
	handler := jsonrpc.NewServer[any, any](ecm, gokitjsonrpctransport.ServerErrorLogger(&logger))
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Post(server.URL, "application/json", addBody())
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d: %s", want, have, buf)
	}
	expectErrorCode(t, gokitjsonrpctransport.InternalError, buf)
	if !logger.Called {
		t.Fatal("Expected logger to be called with error. Wasn't.")
	}
}

func TestServerBadEndpoint(t *testing.T) {
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: func(context.Context, interface{}) (interface{}, error) { return struct{}{}, errors.New("oof") },
			Decode:   nopDecoder,
			Encode:   nopEncoder,
		},
	}
	handler := jsonrpc.NewServer[any, any](ecm)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Post(server.URL, "application/json", addBody())
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	expectErrorCode(t, gokitjsonrpctransport.InternalError, buf)
	expectValidRequestID(t, 1, buf)
}

func TestServerBadEncode(t *testing.T) {
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: gokitendpoint.Nop,
			Decode:   nopDecoder,
			Encode:   func(context.Context, interface{}) (json.RawMessage, error) { return []byte{}, errors.New("oof") },
		},
	}
	handler := jsonrpc.NewServer[any, any](ecm)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Post(server.URL, "application/json", addBody())
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	expectErrorCode(t, gokitjsonrpctransport.InternalError, buf)
	expectValidRequestID(t, 1, buf)
}

func TestServerErrorEncoder(t *testing.T) {
	errTeapot := errors.New("teapot")
	code := func(err error) int {
		if err == errTeapot {
			return http.StatusTeapot
		}
		return http.StatusInternalServerError
	}
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: func(context.Context, interface{}) (interface{}, error) { return struct{}{}, errTeapot },
			Decode:   nopDecoder,
			Encode:   nopEncoder,
		},
	}
	handler := jsonrpc.NewServer[any, any](
		ecm,
		gokitjsonrpctransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) { w.WriteHeader(code(err)) }),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Post(server.URL, "application/json", addBody())
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestCanRejectNonPostRequest(t *testing.T) {
	ecm := gokitjsonrpctransport.EndpointCodecMap{}
	handler := jsonrpc.NewServer[any, any](ecm)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusMethodNotAllowed, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestCanRejectInvalidJSON(t *testing.T) {
	ecm := gokitjsonrpctransport.EndpointCodecMap{}
	handler := jsonrpc.NewServer[any, any](ecm)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Post(server.URL, "application/json", body("clearlynotjson"))
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	expectErrorCode(t, gokitjsonrpctransport.ParseError, buf)
	expectNilRequestID(t, buf)
}

func TestServerUnregisteredMethod(t *testing.T) {
	ecm := gokitjsonrpctransport.EndpointCodecMap{}
	handler := jsonrpc.NewServer[any, any](ecm)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Post(server.URL, "application/json", addBody())
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	expectErrorCode(t, gokitjsonrpctransport.MethodNotFoundError, buf)
}

func TestServerHappyPath(t *testing.T) {
	step, response := testServer(t)
	step()
	resp := <-response
	defer resp.Body.Close() // nolint
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d (%s)", want, have, buf)
	}
	r, err := unmarshalResponse(buf)
	if err != nil {
		t.Fatalf("Can't decode response. err=%s, body=%s", err, buf)
	}
	if r.JSONRPC != gokitjsonrpctransport.Version {
		t.Fatalf("JSONRPC Version: want=%s, got=%s", gokitjsonrpctransport.Version, r.JSONRPC)
	}
	if r.Error != nil {
		t.Fatalf("Unxpected error on response: %s", buf)
	}
}

func TestMultipleServerBeforeCodec(t *testing.T) {
	var done = make(chan struct{})
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: gokitendpoint.Nop,
			Decode:   nopDecoder,
			Encode:   nopEncoder,
		},
	}
	handler := jsonrpc.NewServer[any, any](
		ecm,
		gokitjsonrpctransport.ServerBeforeCodec(func(ctx context.Context, r *http.Request, req gokitjsonrpctransport.Request) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gokitjsonrpctransport.ServerBeforeCodec(func(ctx context.Context, r *http.Request, req gokitjsonrpctransport.Request) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerBeforeCodecs are used")
			}

			close(done)
			return ctx
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	http.Post(server.URL, "application/json", addBody()) // nolint

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestMultipleServerBefore(t *testing.T) {
	var done = make(chan struct{})
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: gokitendpoint.Nop,
			Decode:   nopDecoder,
			Encode:   nopEncoder,
		},
	}
	handler := jsonrpc.NewServer[any, any](
		ecm,
		gokitjsonrpctransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gokitjsonrpctransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerBefores are used")
			}

			close(done)
			return ctx
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	http.Post(server.URL, "application/json", addBody()) // nolint

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestMultipleServerAfter(t *testing.T) {
	var done = make(chan struct{})
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: gokitendpoint.Nop,
			Decode:   nopDecoder,
			Encode:   nopEncoder,
		},
	}
	handler := jsonrpc.NewServer[any, any](
		ecm,
		gokitjsonrpctransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gokitjsonrpctransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerAfters are used")
			}

			close(done)
			return ctx
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	http.Post(server.URL, "application/json", addBody()) // nolint

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestCanFinalize(t *testing.T) {
	var done = make(chan struct{})
	var finalizerCalled bool
	ecm := gokitjsonrpctransport.EndpointCodecMap{
		"add": gokitjsonrpctransport.EndpointCodec{
			Endpoint: gokitendpoint.Nop,
			Decode:   nopDecoder,
			Encode:   nopEncoder,
		},
	}
	handler := jsonrpc.NewServer[any, any](
		ecm,
		gokitjsonrpctransport.ServerFinalizer(func(ctx context.Context, code int, req *http.Request) {
			finalizerCalled = true
			close(done)
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	http.Post(server.URL, "application/json", addBody()) // nolint

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}

	if !finalizerCalled {
		t.Fatal("Finalizer was not called.")
	}
}

func testServer(t *testing.T) (step func(), resp <-chan *http.Response) {
	var (
		stepch   = make(chan bool)
		endpoint = func(ctx context.Context, request interface{}) (response interface{}, err error) {
			<-stepch
			return struct{}{}, nil
		}
		response = make(chan *http.Response)
		ecm      = gokitjsonrpctransport.EndpointCodecMap{
			"add": gokitjsonrpctransport.EndpointCodec{
				Endpoint: endpoint,
				Decode:   nopDecoder,
				Encode:   nopEncoder,
			},
		}
		handler = jsonrpc.NewServer[any, any](ecm)
	)
	go func() {
		server := httptest.NewServer(handler)
		defer server.Close()
		rb := strings.NewReader(`{"jsonrpc": "2.0", "method": "add", "params": [3, 2], "id": 1}`)
		resp, err := http.Post(server.URL, "application/json", rb)
		if err != nil {
			t.Error(err)
			return
		}
		response <- resp
	}()
	return func() { stepch <- true }, response
}
