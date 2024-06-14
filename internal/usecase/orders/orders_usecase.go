package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
)

//go:generate mockgen -package=orders -source=orders_usecase.go -destination=orders_usecase_mock_test.go
type ordersRepository interface {
	InsertOrder(ctx context.Context, order orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
}

type booksRepository interface {
	GetBookByIDs(ctx context.Context, ids []int64) (map[int64]books.Model, error)
}

type usecase struct {
	ordersRepository ordersRepository
	booksRepository  booksRepository
	cfg              *configs.Config
}

func New(ordersRepository ordersRepository, booksRepository booksRepository, cfg *configs.Config) *usecase {
	return &usecase{ordersRepository: ordersRepository, booksRepository: booksRepository, cfg: cfg}
}

func (u *usecase) InsertOrder(ctx context.Context, order orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	bookIDs := make([]int64, 0)
	for _, item := range order.Items {
		bookIDs = append(bookIDs, item.BookID)
	}

	bookMap, err := u.booksRepository.GetBookByIDs(ctx, bookIDs)
	if err != nil {
		return nil, err
	}

	// validate book price and total price
	totalPrice := float64(0)
	for _, item := range order.Items {
		totalPrice += item.Price * float64(item.Quantity)
		if _, ok := bookMap[item.BookID]; !ok {
			return nil, fmt.Errorf("book with id: %d is not found", item.BookID)
		}
		if item.Price != bookMap[item.BookID].Price {
			return nil, fmt.Errorf("book with id: %d has different price", item.BookID)
		}
	}
	if totalPrice != order.TotalAmount {
		return nil, errors.New("total amount is different, please refresh your cart")
	}

	return u.ordersRepository.InsertOrder(ctx, order)
}
