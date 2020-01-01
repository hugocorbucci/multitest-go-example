package storage

import (
	"context"
	"database/sql"
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
