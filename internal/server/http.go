package server

import (
	"database/sql"
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
	r.HandleFunc("/s/{short:[a-f0-9]+}", handler.shortURL).Methods(http.MethodGet)

	return &Server{r}
}

func (h *httpHandler) helloWorld(w http.ResponseWriter, req *http.Request) {
	msg := "Hello, world"
	fmt.Println("Responding request")
	w.Write([]byte(msg))
}

func (h *httpHandler) shortURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	short := mux.Vars(req)["short"]
	if len(short) != 12 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	longURL, err := h.repo.ExpandShortURL(ctx, short)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found\n"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%+v", err)))
		}
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Header().Set("Location", longURL)
}
