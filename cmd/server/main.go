package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hugocorbucci/multitest-go-example/internal/server"
	"github.com/hugocorbucci/multitest-go-example/internal/storage"

	_ "github.com/go-sql-driver/mysql"
)

const (
	port      = "8080"
	sqlDriver = "mysql"

	maxIdleConns    = 10
	maxConnLifetime = 20 * time.Minute
)

func main() {
	ll := log.New(os.Stdout, "HTTP - ", 0)
	addr := net.JoinHostPort("", port)

	dbConn := os.Getenv("DB_CONN")
	if len(dbConn) == 0 {
		dbConn = "root:sekret@tcp(localhost:3306)/multitest"
	}
	db := storage.InitializeStore(context.Background(), sqlDriver, dbConn, ll, maxIdleConns, maxConnLifetime)
	repo := storage.NewSQLStore(db)

	s := server.NewHTTPServer(repo)
	if err := http.ListenAndServe(addr, s); err != nil {
		ll.Fatal("HTTP(s) server failed")
	}
}
