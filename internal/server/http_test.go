package server_test

import (
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

func TestHomeReturnsHelloWorldLocal(t *testing.T) {
	s := server.NewHTTPServer()

	w := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodGet, "/", nil)
	s.ServeHTTP(w, httpReq)

	require.Equal(t, http.StatusOK, w.Code, "expected status code to match for req %+v", httpReq)
	body := string(w.Body.Bytes())
	assert.Equal(t, "Hello, world", body, "expected body to match")
}

func TestHomeReturnsHelloWorld(t *testing.T) {
	baseURL, stop := startTestingHTTPServer(t)
	defer stop()

	httpReq, err := http.NewRequest(http.MethodGet, baseURL+"/", nil)
	require.NoError(t, err, "could not create GET / request")

	resp, err := http.DefaultClient.Do(httpReq)
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
