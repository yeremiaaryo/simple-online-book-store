package books

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"net/http"
	"strconv"
)

//go:generate mockgen -package=books -source=books_handler.go -destination=books_handler_mock_test.go
type booksUsecase interface {
	GetBooks(ctx context.Context, search string, pageSize, pageIndex int) ([]books.Model, error)
}
type Handler struct {
	booksUsecase booksUsecase
}

func New(booksUsecase booksUsecase) *Handler {
	return &Handler{booksUsecase: booksUsecase}
}

func (h *Handler) GetBooks(c echo.Context) error {
	response := books.GetBookListResponse{}
	search := c.QueryParam("search")

	pageIndex, err := strconv.Atoi(c.QueryParam("page_index"))
	if err != nil {
		pageIndex = 1 // default page index is 1 if error
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		pageSize = 10 // default page size is 10 if error
	}

	bookList, err := h.booksUsecase.GetBooks(c.Request().Context(), search, pageSize, pageIndex)
	if err != nil {
		statusCode := http.StatusInternalServerError
		response.Error = err.Error()
		return c.JSON(statusCode, response)
	}
	response.Result = true
	response.Books = bookList
	return c.JSON(http.StatusOK, response)
}
