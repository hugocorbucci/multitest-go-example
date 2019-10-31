package storage

import (
	"database/sql"
)

// sqlStore encapsulates information necessary to query the database
type sqlStore struct {
	db *sql.DB
}

// Repository defines an interface for interacting with the MySQL instances
type Repository interface {
}

// NewSQLStore constructs a new sqlStore
func NewSQLStore(db *sql.DB) Repository {
	return &sqlStore{
		db: db,
	}
}