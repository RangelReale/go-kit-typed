package http_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RangelReale/go-kit-typed/endpoint"
	httptransport "github.com/RangelReale/go-kit-typed/transport/http"
	"github.com/RangelReale/go-kit-typed/util"
	gokitendpoint "github.com/go-kit/kit/endpoint"
	gokithttptransport "github.com/go-kit/kit/transport/http"
)

type serverReq struct {
	req string
}

type serverResp struct {
	resp string
}

func TestServerBadDecode(t *testing.T) {
	handler := httptransport.NewServer[serverReq, serverResp](
		func(context.Context, serverReq) (serverResp, error) { return serverResp{"resp1"}, nil },
		func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, errors.New("dang") },
		func(context.Context, http.ResponseWriter, serverResp) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerBadEndpoint(t *testing.T) {
	handler := httptransport.NewServer[serverReq, serverResp](
		func(context.Context, serverReq) (serverResp, error) { return serverResp{"req1"}, errors.New("dang") },
		func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, nil },
		func(context.Context, http.ResponseWriter, serverResp) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerBadEncode(t *testing.T) {
	handler := httptransport.NewServer[serverReq, serverResp](
		func(context.Context, serverReq) (serverResp, error) { return serverResp{"resp1"}, nil },
		func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, nil },
		func(context.Context, http.ResponseWriter, serverResp) error { return errors.New("dang") },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerNilReq(t *testing.T) {
	var handlerErr error
	handler := httptransport.NewServer[serverReq, serverResp](
		func(context.Context, serverReq) (serverResp, error) { return serverResp{"req1"}, nil },
		httptransport.DecodeRequestFuncAdapter[serverReq](func(ctx context.Context, r *http.Request) (interface{}, error) {
			return nil, nil
		}),
		func(_ context.Context, w http.ResponseWriter, resp serverResp) error {
			w.WriteHeader(http.StatusTeapot)
			return nil
		},
		gokithttptransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			handlerErr = err
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if handlerErr != nil {
		t.Errorf("no error expected, received %v", handlerErr)
	}
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerBadReq(t *testing.T) {
	var handlerErr error
	handler := httptransport.NewServer[serverReq, serverResp](
		func(context.Context, serverReq) (serverResp, error) { return serverResp{"req1"}, nil },
		httptransport.DecodeRequestFuncAdapter[serverReq](func(ctx context.Context, r *http.Request) (interface{}, error) {
			return "bad_type", nil
		}),
		func(_ context.Context, w http.ResponseWriter, resp serverResp) error {
			w.WriteHeader(http.StatusTeapot)
			return nil
		},
		gokithttptransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			handlerErr = err
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if !errors.Is(handlerErr, util.ErrParameterInvalidType) {
		t.Errorf("expected ErrParameterInvalidType, received %v", handlerErr)
	}
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerNilResp(t *testing.T) {
	var handlerErr error
	handler := httptransport.NewServer[serverReq, serverResp](
		endpoint.Adapter[serverReq, serverResp](func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}),
		func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, nil },
		func(_ context.Context, w http.ResponseWriter, resp serverResp) error {
			w.WriteHeader(http.StatusTeapot)
			return nil
		},
		gokithttptransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			handlerErr = err
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if handlerErr != nil {
		t.Errorf("no error expected, received %v", handlerErr)
	}
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerBadResp(t *testing.T) {
	var handlerErr error
	handler := httptransport.NewServer[serverReq, serverResp](
		endpoint.Adapter[serverReq, serverResp](func(ctx context.Context, req interface{}) (interface{}, error) {
			return "bad_type", nil
		}),
		func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, nil },
		func(_ context.Context, w http.ResponseWriter, resp serverResp) error {
			w.WriteHeader(http.StatusTeapot)
			return nil
		},
		gokithttptransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			handlerErr = err
		}),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if !errors.Is(handlerErr, util.ErrParameterInvalidType) {
		t.Errorf("expected ErrParameterInvalidType, received %v", handlerErr)
	}
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerErrorEncoder(t *testing.T) {
	errTeapot := errors.New("teapot")
	code := func(err error) int {
		if errors.Is(err, errTeapot) {
			return http.StatusTeapot
		}
		return http.StatusInternalServerError
	}
	handler := httptransport.NewServer[serverReq, serverResp](
		func(context.Context, serverReq) (serverResp, error) { return serverResp{"resp1"}, errTeapot },
		func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, nil },
		func(context.Context, http.ResponseWriter, serverResp) error { return nil },
		gokithttptransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) { w.WriteHeader(code(err)) }),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerHappyPath(t *testing.T) {
	step, response := testServer(t)
	step()
	resp := <-response
	defer resp.Body.Close()
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d (%s)", want, have, buf)
	}
}

func TestMultipleServerBefore(t *testing.T) {
	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := httptransport.NewServer[serverReq, serverResp](
		endpoint.Adapter[serverReq, serverResp](gokitendpoint.Nop),
		func(context.Context, *http.Request) (serverReq, error) {
			return serverReq{"req1"}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ serverResp) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		gokithttptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gokithttptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerBefores are used")
			}

			close(done)
			return ctx
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestMultipleServerAfter(t *testing.T) {
	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := httptransport.NewServer[serverReq, any](
		endpoint.Adapter[serverReq, any](gokitendpoint.Nop),
		func(context.Context, *http.Request) (serverReq, error) {
			return serverReq{"req1"}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ any) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		gokithttptransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gokithttptransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerAfters are used")
			}

			close(done)
			return ctx
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestServerFinalizer(t *testing.T) {
	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := httptransport.NewServer[serverReq, any](
		endpoint.Adapter[serverReq, any](gokitendpoint.Nop),
		func(context.Context, *http.Request) (serverReq, error) {
			return serverReq{"req1"}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ any) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		gokithttptransport.ServerFinalizer(func(ctx context.Context, code int, _ *http.Request) {
			if want, have := statusCode, code; want != have {
				t.Errorf("StatusCode: want %d, have %d", want, have)
			}

			responseHeader := ctx.Value(gokithttptransport.ContextKeyResponseHeaders).(http.Header)
			if want, have := headerVal, responseHeader.Get(headerKey); want != have {
				t.Errorf("%s: want %q, have %q", headerKey, want, have)
			}

			responseSize := ctx.Value(gokithttptransport.ContextKeyResponseSize).(int64)
			if want, have := int64(len(responseBody)), responseSize; want != have {
				t.Errorf("response size: want %d, have %d", want, have)
			}

			close(done)
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

type enhancedResponse struct {
	Foo string `json:"foo"`
}

func (e enhancedResponse) StatusCode() int      { return http.StatusPaymentRequired }
func (e enhancedResponse) Headers() http.Header { return http.Header{"X-Edward": []string{"Snowden"}} }

func TestEncodeJSONResponse(t *testing.T) {
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return enhancedResponse{Foo: "bar"}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gokithttptransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := http.StatusPaymentRequired, resp.StatusCode; want != have {
		t.Errorf("StatusCode: want %d, have %d", want, have)
	}
	if want, have := "Snowden", resp.Header.Get("X-Edward"); want != have {
		t.Errorf("X-Edward: want %q, have %q", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := `{"foo":"bar"}`, strings.TrimSpace(string(buf)); want != have {
		t.Errorf("Body: want %s, have %s", want, have)
	}
}

type multiHeaderResponse struct{}

func (_ multiHeaderResponse) Headers() http.Header {
	return http.Header{"Vary": []string{"Origin", "User-Agent"}}
}

func TestAddMultipleHeaders(t *testing.T) {
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return multiHeaderResponse{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gokithttptransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	expect := map[string]map[string]struct{}{"Vary": map[string]struct{}{"Origin": struct{}{}, "User-Agent": struct{}{}}}
	for k, vls := range resp.Header {
		for _, v := range vls {
			delete((expect[k]), v)
		}
		if len(expect[k]) != 0 {
			t.Errorf("Header: unexpected header %s: %v", k, expect[k])
		}
	}
}

type multiHeaderResponseError struct {
	multiHeaderResponse
	msg string
}

func (m multiHeaderResponseError) Error() string {
	return m.msg
}

func TestAddMultipleHeadersErrorEncoder(t *testing.T) {
	errStr := "oh no"
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) {
			return nil, multiHeaderResponseError{msg: errStr}
		},
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gokithttptransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	expect := map[string]map[string]struct{}{"Vary": map[string]struct{}{"Origin": struct{}{}, "User-Agent": struct{}{}}}
	for k, vls := range resp.Header {
		for _, v := range vls {
			delete((expect[k]), v)
		}
		if len(expect[k]) != 0 {
			t.Errorf("Header: unexpected header %s: %v", k, expect[k])
		}
	}
	if b, _ := ioutil.ReadAll(resp.Body); errStr != string(b) {
		t.Errorf("ErrorEncoder: got: %q, expected: %q", b, errStr)
	}
}

type noContentResponse struct{}

func (e noContentResponse) StatusCode() int { return http.StatusNoContent }

func TestEncodeNoContent(t *testing.T) {
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return noContentResponse{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gokithttptransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := http.StatusNoContent, resp.StatusCode; want != have {
		t.Errorf("StatusCode: want %d, have %d", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := 0, len(buf); want != have {
		t.Errorf("Body: want no content, have %d bytes", have)
	}
}

type enhancedError struct{}

func (e enhancedError) Error() string                { return "enhanced error" }
func (e enhancedError) StatusCode() int              { return http.StatusTeapot }
func (e enhancedError) MarshalJSON() ([]byte, error) { return []byte(`{"err":"enhanced"}`), nil }
func (e enhancedError) Headers() http.Header         { return http.Header{"X-Enhanced": []string{"1"}} }

func TestEnhancedError(t *testing.T) {
	handler := httptransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return nil, enhancedError{} },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error { return nil },
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("StatusCode: want %d, have %d", want, have)
	}
	if want, have := "1", resp.Header.Get("X-Enhanced"); want != have {
		t.Errorf("X-Enhanced: want %q, have %q", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := `{"err":"enhanced"}`, strings.TrimSpace(string(buf)); want != have {
		t.Errorf("Body: want %s, have %s", want, have)
	}
}

func TestNoOpRequestDecoder(t *testing.T) {
	resw := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error("Failed to create request")
	}
	handler := httptransport.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			if request != nil {
				t.Error("Expected nil request in hendpoint when using NopRequestDecoder")
			}
			return nil, nil
		},
		gokithttptransport.NopRequestDecoder,
		gokithttptransport.EncodeJSONResponse,
	)
	handler.ServeHTTP(resw, req)
	if resw.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resw.Code)
	}
}

func testServer(t *testing.T) (step func(), resp <-chan *http.Response) {
	var (
		stepch   = make(chan bool)
		endpoint = func(context.Context, serverReq) (serverResp, error) { <-stepch; return serverResp{"resp1"}, nil }
		response = make(chan *http.Response)
		handler  = httptransport.NewServer[serverReq, serverResp](
			endpoint,
			func(context.Context, *http.Request) (serverReq, error) { return serverReq{"req1"}, nil },
			func(context.Context, http.ResponseWriter, serverResp) error { return nil },
			gokithttptransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context { return ctx }),
			gokithttptransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context { return ctx }),
		)
	)
	go func() {
		server := httptest.NewServer(handler)
		defer server.Close()
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Error(err)
			return
		}
		response <- resp
	}()
	return func() { stepch <- true }, response
}
