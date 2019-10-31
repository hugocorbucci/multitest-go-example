package server

import (
	"fmt"
	"net/http"

	"github.com/hugocorbucci/multitest-go-example/internal/storage"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
  *mux.Router
}

type httpHandler struct {
	repo storage.Repository
}

// NewHTTPServer creates a new server
func NewHTTPServer(repo storage.Repository) *Server {
	handler := &httpHandler{repo}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.helloWorld).Methods(http.MethodGet)

	return &Server{r}
}

func (h *httpHandler) helloWorld(w http.ResponseWriter, req *http.Request) {
	msg := "Hello, world"
	fmt.Println("Responding request")	
  w.Write([]byte(msg))
}
