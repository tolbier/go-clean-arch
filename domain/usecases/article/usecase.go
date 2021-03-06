package article

import (
    "context"
	"github.com/tolbier/go-clean-arch/domain"
	"github.com/tolbier/go-clean-arch/domain/entities"
	"github.com/tolbier/go-clean-arch/domain/repositories"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Usecase represent the article's usecases
type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]entities.Article, string, error)
	GetByID(ctx context.Context, id int64) (entities.Article, error)
	Update(ctx context.Context, ar *entities.Article) error
	GetByTitle(ctx context.Context, title string) (entities.Article, error)
	Store(context.Context, *entities.Article) error
	Delete(ctx context.Context, id int64) error
}

type usecase struct {
	articleRepo    repositories.ArticleRepository
	authorRepo     repositories.AuthorRepository
	contextTimeout time.Duration
}

// NewUsecase will create new an usecase object representation of domain.Usecase interface
func NewUsecase(a repositories.ArticleRepository, ar repositories.AuthorRepository, timeout time.Duration) Usecase {
	return &usecase{
		articleRepo:    a,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *usecase) fillAuthorDetails(c context.Context, data []entities.Article) ([]entities.Article, error) {
	g, ctx := errgroup.WithContext(c)

	// Get the author's id
	mapAuthors := map[int64]entities.Author{}

	for _, article := range data {
		mapAuthors[article.Author.ID] = entities.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan entities.Author)
	for authorID := range mapAuthors {
		authorID := authorID
		g.Go(func() error {
			res, err := a.authorRepo.GetByID(ctx, authorID)
			if err != nil {
				return err
			}
			chanAuthor <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanAuthor)
	}()

	for author := range chanAuthor {
		if author != (entities.Author{}) {
			mapAuthors[author.ID] = author
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

func (a *usecase) Fetch(c context.Context, cursor string, num int64) (res []entities.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		nextCursor = ""
	}
	return
}

func (a *usecase) GetByID(c context.Context, id int64) (res entities.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return entities.Article{}, err
	}
	res.Author = resAuthor
	return
}

func (a *usecase) Update(c context.Context, ar *entities.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *usecase) GetByTitle(c context.Context, title string) (res entities.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return entities.Article{}, err
	}

	res.Author = resAuthor
	return
}

func (a *usecase) Store(c context.Context, m *entities.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, m.Title)
	if existedArticle != (entities.Article{}) {
		return domain.ErrConflict
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *usecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (entities.Article{}) {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
