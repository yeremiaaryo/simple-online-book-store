package books

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/internal/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetBooks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBooksUC := NewMockbooksUsecase(mockCtrl)

	// Setup Echo framework context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type args struct {
		search    string
		pageIndex string
		pageSize  string
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		expectedResult books.GetBookListResponse
		mockFn         func(args args)
	}{
		{
			name: "Valid parameters",
			args: args{
				search:    "Harry Potter",
				pageIndex: "1",
				pageSize:  "10",
			},
			expectedStatus: http.StatusOK,
			expectedResult: books.GetBookListResponse{
				BaseResponse: response.BaseResponse{
					Result: true,
					Error:  "",
				},
				Books: []books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789"},
				},
			},
			mockFn: func(args args) {
				mockBooksUC.EXPECT().GetBooks(c.Request().Context(), args.search, 10, 1).Return([]books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				}, nil)
			},
		},
		{
			name: "Invalid pageIndex (not a number), continue with default value",
			args: args{
				search:    "Harry Potter",
				pageIndex: "invalid",
				pageSize:  "10",
			},
			expectedStatus: http.StatusOK,
			expectedResult: books.GetBookListResponse{
				BaseResponse: response.BaseResponse{
					Result: true,
					Error:  "",
				},
				Books: []books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789"},
				},
			},
			mockFn: func(args args) {
				mockBooksUC.EXPECT().GetBooks(c.Request().Context(), args.search, 10, 1).Return([]books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				}, nil)
			},
		},
		{
			name: "Invalid pageSize (not a number)",
			args: args{
				search:    "Harry Potter",
				pageIndex: "1",
				pageSize:  "invalid",
			},
			expectedStatus: http.StatusOK,
			expectedResult: books.GetBookListResponse{
				BaseResponse: response.BaseResponse{
					Result: true,
					Error:  "",
				},
				Books: []books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789"},
				},
			},
			mockFn: func(args args) {
				mockBooksUC.EXPECT().GetBooks(c.Request().Context(), args.search, 10, 1).Return([]books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				}, nil)
			},
		},
		{
			name: "Error from usecase",
			args: args{
				search:    "",
				pageIndex: "1",
				pageSize:  "10",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResult: books.GetBookListResponse{
				BaseResponse: response.BaseResponse{
					Result: false,
					Error:  "mock error from usecase",
				},
			},
			mockFn: func(args args) {
				mockBooksUC.EXPECT().GetBooks(c.Request().Context(), args.search, 10, 1).Return(nil, errors.New("mock error from usecase"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set query parameters in the request
			q := req.URL.Query()
			q.Set("search", tt.args.search)
			q.Set("page_index", tt.args.pageIndex)
			q.Set("page_size", tt.args.pageSize)
			req.URL.RawQuery = q.Encode()

			// Reset recorder for each test case
			rec = httptest.NewRecorder()
			c = e.NewContext(req, rec)

			tt.mockFn(tt.args)
			h := &Handler{
				booksUsecase: mockBooksUC,
			}
			err := h.GetBooks(c)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, rec.Code)
			var result books.GetBookListResponse
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult.Result, result.Result)
			assert.Equal(t, tt.expectedResult.Books, result.Books)

		})
	}
}
