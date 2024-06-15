package books

import (
	"context"
	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/pkg/util"
)

//go:generate mockgen -package=books -source=books_usecase.go -destination=books_usecase_mock_test.go
type booksRepository interface {
	GetBooks(ctx context.Context, search string, limit, offset int) ([]books.Model, error)
}

type usecase struct {
	booksRepository booksRepository
	cfg             *configs.Config
}

func New(booksRepository booksRepository, cfg *configs.Config) *usecase {
	return &usecase{booksRepository: booksRepository, cfg: cfg}
}

func (u *usecase) GetBooks(ctx context.Context, search string, pageSize, pageIndex int) ([]books.Model, error) {
	// convert to limit and offset
	limit, offset := util.GetLimitAndOffset(pageIndex, pageSize)
	return u.booksRepository.GetBooks(ctx, search, limit, offset)
}
