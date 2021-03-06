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
	"time"

	"github.com/hugocorbucci/multitest-go-example/internal/server"
	"github.com/hugocorbucci/multitest-go-example/internal/storage"
	"github.com/hugocorbucci/multitest-go-example/internal/storage/sqltesting"
	"github.com/hugocorbucci/multitest-go-example/internal/storage/stubs"

	_ "github.com/go-sql-driver/mysql"
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
		input := "123456789012"
		output := "https://www.digitalocean.com"
		do(func() {
			err := deps.DB.RegisterURLMapping(context.Background(), output, input)
			require.NoError(t, err, "unexpected error registering url")

			httpReq, err := http.NewRequest(http.MethodGet, deps.BaseURL+"/s/"+input, nil)
			require.NoError(t, err, "could not create GET / request")
			resp, err := deps.HTTPClient.Do(httpReq)
			require.NoError(t, err, "error making request %+v", httpReq)
			require.Equal(t, http.StatusFound, resp.StatusCode, "expected status code to match for req %+v", httpReq)

			body, err := readBodyFrom(resp)
			require.NoError(t, err, "unexpected error reading response body")
			assert.Equal(t, "", body, "expected body to match")
			assert.Equal(t, output, resp.Header.Get("Location"), "expected location to match")
		}).withTearDown(func() {
			err := deps.DB.ClearMappingWithKey(context.Background(), input)
			require.NoError(t, err, "unexpected error clearing mapping")
		}).Now()
	})
}

// TestDependencies encapsulates the dependencies needed to run a test
type TestDependencies struct {
	BaseURL    string
	HTTPClient HTTPClient
	DB         storage.Repository
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

type testStructure struct {
	test     func()
	tearDown func()
}

func do(test func()) *testStructure {
	return &testStructure{
		test: test,
	}
}

func withTearDown(tearDown func()) *testStructure {
	return &testStructure{
		tearDown: tearDown,
	}
}

func (s *testStructure) do(test func()) *testStructure {
	copy := &testStructure{}
	*copy = *s
	copy.test = test
	return copy
}

func (s *testStructure) withTearDown(tearDown func()) *testStructure {
	copy := &testStructure{}
	*copy = *s
	copy.tearDown = tearDown
	return copy
}

func (s *testStructure) Now() {
	if s.tearDown != nil {
		defer s.tearDown()
	}

	if s.test != nil {
		s.test()
	}
}

func unitDependencies(*testing.T) (*TestDependencies, func()) {
	repo := stubs.NewStubRepository(nil)
	s := server.NewHTTPServer(repo)
	httpClient := &InMemoryHTTPClient{server: s}
	return &TestDependencies{
		BaseURL:    "",
		HTTPClient: httpClient,
		DB:         repo,
	}, func() {}
}

func integrationDependencies(t *testing.T) (*TestDependencies, func()) {
	ctx := context.Background()
	inMemoryStore, err := sqltesting.NewMemoryStore(ctx, sqltesting.NewTestingLog(t))
	require.NoError(t, err, "error creating in memory store")

	baseURL, stop := startTestingHTTPServer(t, inMemoryStore)
	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestDependencies{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		DB:         inMemoryStore.Store,
	}, stop
}

func smokeDependencies(t *testing.T) *TestDependencies {
	ctx := context.Background()
	l := sqltesting.NewTestingLog(t)
	db := storage.InitializeStore(ctx, "mysql", os.Getenv("DB_CONN"), l, 1, 10*time.Second)

	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &TestDependencies{
		BaseURL:    os.Getenv("TARGET_URL"),
		HTTPClient: http.DefaultClient,
		DB:         storage.NewSQLStore(db),
	}
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
