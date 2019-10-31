package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hugocorbucci/multitest-go-example/internal/server"
	"github.com/hugocorbucci/multitest-go-example/internal/storage"

	_ "github.com/go-sql-driver/mysql"
)

const (
	port      = "8080"
	sqlDriver = "mysql"
)

func main() {
	ll := log.New(os.Stdout, "HTTP - ", 0)
	addr := net.JoinHostPort("", port)

	dbConn := os.Getenv("DB_CONN")
	if len(dbConn) == 0 {
		dbConn = "root:sekret@tcp(localhost:3306)/multitest"
	}
	db := initializeStore(context.Background(), dbConn, ll)
	repo := storage.NewSQLStore(db)

	s := server.NewHTTPServer(repo)
	if err := http.ListenAndServe(addr, s); err != nil {
		ll.Fatal("HTTP(s) server failed")
	}
}

func initializeStore(ctx context.Context, conn string, l *log.Logger) *sql.DB {
	mainDB, err := sql.Open(sqlDriver, conn)
	if err != nil {
		l.Fatalf("failed to open main store: %v", err)
	}

	row := mainDB.QueryRowContext(ctx, "SELECT 1 FROM dual")
	var ping int
	err = row.Scan(&ping)
	if err != nil {
		l.Fatalf("failed to ping DB: %v", err)
	}

	return mainDB
}
