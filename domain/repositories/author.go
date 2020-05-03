package repositories

import (
	"context"
	"github.com/tolbier/go-clean-arch/domain/entities"
)

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (entities.Author, error)
}
