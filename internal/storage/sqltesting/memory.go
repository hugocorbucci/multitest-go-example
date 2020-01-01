package sqltesting

import (
	"context"
	"database/sql"
	"io"
	"log"
	"testing"
	"time"

	"github.com/hugocorbucci/multitest-go-example/internal/storage"

	// Import sqlite only for memory usage
	_ "github.com/mattn/go-sqlite3"
)

var _ storage.Repository = storage.NewSQLStore(nil)

// MemoryStore is a SQLite in-memory representation of a Repository.
type MemoryStore struct {
	DB    *sql.DB
	Store storage.Repository
}

// NewMemoryStore returns a sqlite "in-memory" repository for local tests.
func NewMemoryStore(ctx context.Context, l *log.Logger) (*MemoryStore, error) {
	db, err := NewSQLiteDB(ctx, l)
	if err != nil {
		return nil, err
	}

	return &MemoryStore{
		DB:    db,
		Store: storage.NewSQLStore(db),
	}, nil
}

// NewSQLiteDB creates a SQLite DB with the default structure
func NewSQLiteDB(ctx context.Context, l *log.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, MemoryStructure)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type testingLogWriter struct {
	t *testing.T
}

var _ io.Writer = &testingLogWriter{}

func (w *testingLogWriter) Write(b []byte) (int, error) {
	w.t.Log(string(b))
	return len(b), nil
}

// NewTestingLog creates a logger that outputs to the testing framework
func NewTestingLog(t *testing.T) *log.Logger {
	return log.New(&testingLogWriter{t}, "", 0)
}

// MemoryStructure is the database structure for a the sqlite "memory" database.
const MemoryStructure = `
PRAGMA synchronous = OFF;
PRAGMA journal_mode = OFF;

CREATE TABLE IF NOT EXISTS url_mapping (
	id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	url VARCHAR(255) NOT NULL,
	short_url VARCHAR(12) NOT NULL,
	updated_at DEFAULT CURRENT_TIMESTAMP,
	created_at DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_url_mapping_short_url" ON "url_mapping" (short_url);
`
