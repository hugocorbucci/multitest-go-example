package server_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hugocorbucci/multitest-go-example/internal/server"
	"github.com/hugocorbucci/multitest-go-example/internal/storage/sqltesting"
	"github.com/hugocorbucci/multitest-go-example/internal/storage/stubs"

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

func TestHomeReturnsHelloWorld(t *testing.T) {
	withDependencies(t, func(t *testing.T, deps *TestDependencies) {
		httpReq, err := http.NewRequest(http.MethodGet, deps.BaseURL+"/", nil)
		require.NoError(t, err, "could not create GET / request")

		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)

		require.Equal(t, http.StatusOK, resp.StatusCode, "expected status code to match for req %+v", httpReq)
		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "Hello, world", body, "expected body to match")
	})
}

func TestShortURLReturnsNotFoundForInvalidURL(baseT *testing.T) {
	withDependencies(baseT, func(t *testing.T, deps *TestDependencies) {
		httpReq, err := http.NewRequest(http.MethodGet, deps.BaseURL+"/s/invalid", nil)
		require.NoError(t, err, "could not create GET / request")
		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "expected status code to match for req %+v", httpReq)

		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "404 page not found\n", body, "expected body to match")
	})
}

func TestShortURLReturnsNotFoundForUnknown(baseT *testing.T) {
	withDependencies(baseT, func(t *testing.T, deps *TestDependencies) {
		httpReq, err := http.NewRequest(http.MethodGet, deps.BaseURL+"/s/123456789012", nil)
		require.NoError(t, err, "could not create GET / request")
		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "expected status code to match for req %+v", httpReq)

		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "404 page not found\n", body, "expected body to match")
	})
}

func TestShortURLReturnsFoundForValidURL(baseT *testing.T) {
	withDependencies(baseT, func(t *testing.T, deps *TestDependencies) {
		mocking := false
		input := "123456789012"
		output := "https://www.digitalocean.com"

		if mocking {
			// TODO: Mockar
		} else {
			// TODO: cadastrar a URL
		}

		httpReq, err := http.NewRequest(http.MethodGet, deps.BaseURL+"/s/"+input, nil)
		require.NoError(t, err, "could not create GET / request")
		resp, err := deps.HTTPClient.Do(httpReq)
		require.NoError(t, err, "error making request %+v", httpReq)
		require.Equal(t, http.StatusFound, resp.StatusCode, "expected status code to match for req %+v", httpReq)

		body, err := readBodyFrom(resp)
		require.NoError(t, err, "unexpected error reading response body")
		assert.Equal(t, "", body, "expected body to match")
		assert.Equal(t, output, resp.Header.Get("Location"), "expected location to match")
	})
}

// TestDependencies encapsulates the dependencies needed to run a test
type TestDependencies struct {
	BaseURL    string
	HTTPClient HTTPClient
}

func withDependencies(baseT *testing.T, test func(*testing.T, *TestDependencies)) {
	if len(os.Getenv("TARGET_URL")) == 0 {
		testStates := map[string]func(*testing.T) (*TestDependencies, func()){
			"unitServerTest":        unitDependencies,
			"integrationServerTest": integrationDependencies,
		}
		for name, dep := range testStates {
			baseT.Run(name, func(t *testing.T) {
				deps, stop := dep(t)
				defer stop()
				test(t, deps)
			})
		}
	} else {
		test(baseT, smokeDependencies(baseT))
	}
}

func unitDependencies(*testing.T) (*TestDependencies, func()) {
	repo := stubs.NewStubRepository(nil)
	s := server.NewHTTPServer(repo)
	httpClient := &InMemoryHTTPClient{server: s}
	return &TestDependencies{BaseURL: "", HTTPClient: httpClient}, func() {}
}

func integrationDependencies(t *testing.T) (*TestDependencies, func()) {
	ctx := context.Background()
	inMemoryStore, err := sqltesting.NewMemoryStore(ctx, sqltesting.NewTestingLog(t))
	require.NoError(t, err, "error creating in memory store")
	baseURL, stop := startTestingHTTPServer(t, inMemoryStore)
	return &TestDependencies{BaseURL: baseURL, HTTPClient: http.DefaultClient}, stop
}

func smokeDependencies(_ *testing.T) *TestDependencies {
	return &TestDependencies{BaseURL: os.Getenv("TARGET_URL"), HTTPClient: http.DefaultClient}
}

func startTestingHTTPServer(t *testing.T, inMemoryStore *sqltesting.MemoryStore) (string, func()) {
	ctx := context.Background()
	s := server.NewHTTPServer(inMemoryStore.Store)

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
