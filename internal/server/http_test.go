package server_test

import (
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