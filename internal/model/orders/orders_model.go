package orders

import "github.com/yeremiaaryo/gotu-assignment/internal/response"

type OrderStatus string

const (
	OrderStatusNew OrderStatus = "NEW"
)

func (os OrderStatus) String() string {
	return string(os)
}

type (
	Model struct {
		ID          int64   `db:"id"`
		UserID      int64   `db:"user_id"`
		TotalAmount float64 `db:"total_amount"`
		Status      string  `db:"status"`
		CreatedAt   int64   `db:"created_at"`
		UpdatedAt   int64   `db:"updated_at"`
	}

	OrderItem struct {
		ID        int64   `db:"id"`
		OrderID   int64   `db:"order_id"`
		BookID    int64   `db:"book_id"`
		Quantity  int     `db:"quantity"`
		Price     float64 `db:"price"`
		CreatedAt int64   `db:"created_at"`
		UpdatedAt int64   `db:"updated_at"`
	}

	History struct {
		ID          int64         `json:"order_id"`
		TotalAmount float64       `json:"total_amount"`
		Status      string        `json:"status"`
		CreatedAt   int64         `json:"created_at"`
		UpdatedAt   int64         `json:"updated_at"`
		Items       []ItemHistory `json:"items"`
	}

	ItemHistory struct {
		ID       int64   `json:"item_id"`
		BookID   int64   `json:"book_id"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
)

type (
	CreateOrderRequest struct {
		UserID      int64             `json:"-"`
		TotalAmount float64           `json:"total_amount" validate:"required"`
		Items       []CreateOrderItem `json:"items" validate:"required"`
	}

	CreateOrderItem struct {
		BookID   int64   `json:"book_id" validate:"required"`
		Quantity int     `json:"quantity" validate:"required"`
		Price    float64 `json:"price" validate:"required"`
	}
)

type (
	CreateOrderResponse struct {
		response.BaseResponse
		OrderID int64  `json:"order_id"`
		Status  string `json:"status"`
	}

	OrderHistoryResponse struct {
		response.BaseResponse
		Histories []History `json:"data"`
	}
)
