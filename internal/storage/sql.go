package storage

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// sqlStore encapsulates information necessary to query the database
type sqlStore struct {
	db *sql.DB
}

// Repository defines an interface for interacting with the MySQL instances
type Repository interface {
	// ExpandShortURL returns the stored full URL value for the given short path
	ExpandShortURL(context.Context, string) (string, error)
}

// NewSQLStore constructs a new sqlStore
func NewSQLStore(db *sql.DB) Repository {
	return &sqlStore{
		db: db,
	}
}

func (s *sqlStore) ExpandShortURL(_ context.Context, _ string) (string, error) {
	return "", sql.ErrNoRows
}

// InitializeStore initializes a sql connection to a given connection
func InitializeStore(ctx context.Context, sqlDriver, conn string, l *log.Logger, maxIdleConns int, maxConnLifetime time.Duration) *sql.DB {
	mainDB, err := sql.Open(sqlDriver, conn)
	if err != nil {
		l.Fatalf("failed to open main store: %v", err)
	}
	mainDB.SetMaxIdleConns(maxIdleConns)
	mainDB.SetConnMaxLifetime(maxConnLifetime)

	row := mainDB.QueryRowContext(ctx, "SELECT 1 FROM dual")
	var ping int
	err = row.Scan(&ping)
	if err != nil {
		l.Fatalf("failed to ping DB: %v", err)
	}

	return mainDB
}
