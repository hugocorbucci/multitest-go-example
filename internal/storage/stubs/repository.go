package stubs

import (
	"context"
	"database/sql"

	"github.com/hugocorbucci/multitest-go-example/internal/storage"
)

var _ storage.Repository = &Repository{}

// NewStubRepository creates a new repository with base data
func NewStubRepository(initialData map[string]string) *Repository {
	data := initialData
	if data == nil {
		data = make(map[string]string)
	}
	return &Repository{
		data: data,
	}
}

// Repository is a repository stub
type Repository struct {
	data map[string]string
}

// ExpandShortURL returns the long string at the given key or a sql.NoRowsError if none is found
func (r *Repository) ExpandShortURL(_ context.Context, s string) (string, error) {
	v, ok := r.data[s]
	if !ok {
		return "", sql.ErrNoRows
	}
	return v, nil
}

// Add sets the input as the key and the output and the value for that key
func (r *Repository) Add(input, output string) {
	r.data[input] = output
}

// RegisterURLMapping adds the shortPath as the key for the longURL
func (r *Repository) RegisterURLMapping(_ context.Context, longURL, shortPath string) error {
	r.Add(shortPath, longURL)
	return nil
}

// ClearMappingWithKey removes any entry associated with the given short path
func (r *Repository) ClearMappingWithKey(_ context.Context, shortPath string) error {
	delete(r.data, shortPath)
	return nil
}
