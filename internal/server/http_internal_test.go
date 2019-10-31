package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWord(t *testing.T) {
	handler := &httpHandler{}

	w := httptest.NewRecorder()

	handler.helloWorld(w, nil)

	require.Equal(t, http.StatusOK, w.Code, "expected status code to match")
	body := string(w.Body.Bytes())
	assert.Equal(t, "Hello, world", body, "expected body to match")
}