package article_test

import (
    "context"
    "errors"
    "github.com/tolbier/go-clean-arch/domain"
    "github.com/tolbier/go-clean-arch/domain/entities"
    "github.com/tolbier/go-clean-arch/domain/usecases/article"
    . "github.com/tolbier/go-clean-arch/mocks/domain/repositories"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockArticleRepo := new(ArticleRepository)
	mockArticle := entities.Article{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]entities.Article, 0)
	mockListArtilce = append(mockListArtilce, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		mockAuthor := entities.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(ArticleRepository)
	mockArticle := entities.Article{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := entities.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()
		mockAuthorrepo := new(AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entities.Article{}, errors.New("Unexpected")).Once()

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, entities.Article{}, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockArticleRepo := new(ArticleRepository)
	mockArticle := entities.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(entities.Article{}, domain.ErrNotFound).Once()
		mockArticleRepo.On("Store", mock.Anything, mock.AnythingOfType("*entities.Article")).Return(nil).Once()

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()
		mockAuthor := entities.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(ArticleRepository)
	mockArticle := entities.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()

		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("article-is-not-exist", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entities.Article{}, nil).Once()

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(entities.Article{}, errors.New("Unexpected Error")).Once()

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockArticleRepo := new(ArticleRepository)
	mockArticle := entities.Article{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("Update", mock.Anything, &mockArticle).Once().Return(nil)

		mockAuthorrepo := new(AuthorRepository)
		u := article.NewUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.Update(context.TODO(), &mockArticle)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
