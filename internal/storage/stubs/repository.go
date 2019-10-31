package stubs

import (
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
