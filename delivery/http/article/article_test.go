package article_test

import (
    "encoding/json"
    "github.com/tolbier/go-clean-arch/delivery/http/article"
    "github.com/tolbier/go-clean-arch/domain"
    "github.com/tolbier/go-clean-arch/domain/entities"
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"
    "time"

    "github.com/bxcodec/faker"
    "github.com/labstack/echo"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    . "github.com/tolbier/go-clean-arch/mocks/domain/usecases/article"
)

func TestFetch(t *testing.T) {
	var mockArticle entities.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockUCase := new(Usecase)
	mockListArticle := make([]entities.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListArticle, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := article.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(Usecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := article.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchArticle(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockArticle entities.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(Usecase)

	num := int(mockArticle.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockArticle, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := article.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockArticle := entities.Article{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	mockUCase := new(Usecase)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*entities.Article")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/article", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/article")

	handler := article.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockArticle entities.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)

	mockUCase := new(Usecase)

	num := int(mockArticle.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/article/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("article/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := article.ArticleHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
