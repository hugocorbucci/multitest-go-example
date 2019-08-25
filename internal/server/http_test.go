package server_test

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hugocorbucci/multitest-go-example/internal/server"
	"github.com/stretchr/testify/require"
)

var _ http.Handler = server.NewHTTPServer()

func TestHomeReturnsHelloWorld(t *testing.T) {
	addr, stop := startTestingHTTPServer(t)
	defer stop()

	httpReq, err := http.NewRequest(http.MethodGet, addr+"/", nil)
	require.NoError(t, err, "could not create GET / request")
	resp, err := http.DefaultClient.Do(httpReq)

	require.NoError(t, err, "error making request %+v", httpReq)
	require.Equal(t, 200, resp.StatusCode, "expected status code to match for req %+v", httpReq)
	body, err := readBodyFrom(resp)
	require.NoError(t, err, "unexpected error reading response body")
	assert.Equal(t, "Hello, world", body, "expected body to match")
}

func startTestingHTTPServer(t *testing.T) (string, func()) {
	s := server.NewHTTPServer()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("could not listen for HTTP requests: %+v", err)
	}
	addr := "http://" + listener.Addr().String()
	srvr := http.Server{Addr: addr, Handler: s}

	go srvr.Serve(listener)
	return addr, func() {
		if err := srvr.Shutdown(context.Background()); err != nil {
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
