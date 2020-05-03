package repositories

import (
	"context"
	. "github.com/tolbier/go-clean-arch/domain/entities"
)

// ArticleRepository represent the article's repository contract
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Article, error)
	GetByTitle(ctx context.Context, title string) (Article, error)
	Update(ctx context.Context, ar *Article) error
	Store(ctx context.Context, a *Article) error
	Delete(ctx context.Context, id int64) error
}
