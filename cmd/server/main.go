package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hugocorbucci/multitest-go-example/internal/server"
)

const (
	port = "8080"
)

func main() {
	s := server.NewHTTPServer()
	ll := log.New(os.Stdout, "HTTP", 0)
	addr := net.JoinHostPort("", port)

	if err := http.ListenAndServe(addr, s); err != nil {
		ll.Fatal("HTTP(s) server failed")
	}
}
