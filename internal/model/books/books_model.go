package books

import (
	"github.com/yeremiaaryo/gotu-assignment/internal/response"
	"time"
)

type (
	Model struct {
		ID            int64     `json:"id" db:"id"`
		Title         string    `json:"title" db:"title"`
		Author        string    `json:"author" db:"author"`
		ISBN          string    `json:"isbn" db:"isbn"`
		PublishedDate time.Time `json:"published_date" db:"published_date"`
		Price         float64   `json:"price" db:"price"`
		CreatedAt     int64     `json:"-" db:"created_at"`
		UpdatedAt     int64     `json:"-" db:"updated_at"`
	}
)

type (
	GetBookListResponse struct {
		response.BaseResponse
		Books []Model `json:"books"`
	}
)
