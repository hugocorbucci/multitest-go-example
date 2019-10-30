package server_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hugocorbucci/multitest-go-example/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}
type InMemoryHTTPClient struct {
	server *server.Server
}

func (c *InMemoryHTTPClient) Do(r *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	c.server.ServeHTTP(w, r)

	return &http.Response{
		Header:     w.Header(),
		StatusCode: w.Code,
		Body:       ioutil.NopCloser(bytes.NewReader(w.Body.Bytes())),
	}, nil
}

func TestHomeReturnsHelloWorldLocal(t *testing.T) {
	s := server.NewHTTPServer()
	httpClient := &InMemoryHTTPClient{server: s}
	baseURL := ""

	httpReq, err := http.NewRequest(http.MethodGet, baseURL+"/", nil)
	require.NoError(t, err, "could not create GET / request")

	resp, err := httpClient.Do(httpReq)
	require.NoError(t, err, "error making request %+v", httpReq)
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected status code to match for req %+v", httpReq)
	body, err := readBodyFrom(resp)
	require.NoError(t, err, "unexpected error reading response body")
	assert.Equal(t, "Hello, world", body, "expected body to match")
}

func TestHomeReturnsHelloWorld(t *testing.T) {
	baseURL, stop := startTestingHTTPServer(t)
	defer stop()
	httpClient := http.DefaultClient

	httpReq, err := http.NewRequest(http.MethodGet, baseURL+"/", nil)
	require.NoError(t, err, "could not create GET / request")

	resp, err := httpClient.Do(httpReq)
	require.NoError(t, err, "error making request %+v", httpReq)
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected status code to match for req %+v", httpReq)
	body, err := readBodyFrom(resp)
	require.NoError(t, err, "unexpected error reading response body")
	assert.Equal(t, "Hello, world", body, "expected body to match")
}

func startTestingHTTPServer(t *testing.T) (string, func()) {
	ctx := context.Background()
	s := server.NewHTTPServer()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("could not listen for HTTP requests: %+v", err)
	}
	baseURL := "http://" + listener.Addr().String()
	srvr := http.Server{Addr: baseURL, Handler: s}

	go srvr.Serve(listener)
	return baseURL, func() {
		if err := srvr.Shutdown(ctx); err != nil {
			t.Logf("could not shutdown http server: %+v", err)
		}
	}
}


func readBodyFrom(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
