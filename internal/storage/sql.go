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
	// RegisterURLMapping registers the first string as the long value URL on the second string short path.
	// Emits an error if parameters are invalid or registration fails
	RegisterURLMapping(context.Context, string, string) error
}

// NewSQLStore constructs a new sqlStore
func NewSQLStore(db *sql.DB) Repository {
	return &sqlStore{
		db: db,
	}
}

// ExpandShortURL returns the corresponding long url for a short path or
// sql.ErrNoRows if no matching path is found
func (s *sqlStore) ExpandShortURL(ctx context.Context, shortPath string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT url FROM url_mapping WHERE short_url = ?", shortPath)
	var result string
	err := row.Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// RegisterURLMapping associates a short path with a long url
// It returns an error if the short URL already exists
func (s *sqlStore) RegisterURLMapping(ctx context.Context, longURL, shortPath string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO url_mapping (url, short_url) VALUES (?, ?)", longURL, shortPath)
	if err != nil {
		return err
	}
	return nil
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
