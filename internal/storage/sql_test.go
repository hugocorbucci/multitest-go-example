package storage_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hugocorbucci/multitest-go-example/internal/storage"
	"github.com/hugocorbucci/multitest-go-example/internal/storage/sqltesting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandShortURLWithoutMatch(t *testing.T) {
	ctx := context.Background()
	db, err := sqltesting.NewSQLiteDB(ctx, sqltesting.NewTestingLog(t))
	require.NoError(t, err, "error building db")
	r := storage.NewSQLStore(db)
	_, err = r.ExpandShortURL(ctx, "")
	assert.EqualError(t, err, sql.ErrNoRows.Error(), "expecte error to match")
}

func TestExpandShortURLWithMatch(t *testing.T) {
	ctx := context.Background()
	db, err := sqltesting.NewSQLiteDB(ctx, sqltesting.NewTestingLog(t))
	require.NoError(t, err, "error building db")
	input := "123456789012"
	output := "https://www.digitalocean.com"
	r, err := db.ExecContext(ctx, "INSERT INTO url_mapping (url, short_url) VALUES (?, ?)", output, input)
	require.NoError(t, err, "error inserting data")
	ra, err := r.RowsAffected()
	require.NoError(t, err, "error getting rows affected")
	assert.EqualValues(t, 1, ra, "expected rows affected to match")
	store := storage.NewSQLStore(db)
	s, err := store.ExpandShortURL(ctx, input)
	require.NoError(t, err, "expected no error")
	assert.Equal(t, output, s, "expected short url to match")
}
