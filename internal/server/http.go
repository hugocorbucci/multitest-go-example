package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
  *mux.Router
}
	
// NewHTTPServer creates a new server
func NewHTTPServer() *Server {
  r := mux.NewRouter()
  r.HandleFunc("/", helloWorld).Methods(http.MethodGet)

  return &Server{r}
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	msg := "Hello, world"
	fmt.Println("Responding request")	
  w.Write([]byte(msg))
}
