package orders

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
	"github.com/yeremiaaryo/gotu-assignment/pkg/util"
	"net/http"
)

//go:generate mockgen -package=orders -source=orders_handler.go -destination=orders_handler_mock_test.go
type ordersUsecase interface {
	InsertOrder(ctx context.Context, order orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
}
type Handler struct {
	ordersUsecase ordersUsecase
}

func New(ordersUsecase ordersUsecase) *Handler {
	return &Handler{ordersUsecase: ordersUsecase}
}

func (h *Handler) CreateOrder(c echo.Context) error {
	response := orders.CreateOrderResponse{}

	userID, err := util.GetUserID(c)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	var request orders.CreateOrderRequest
	err = c.Bind(&request)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	request.UserID = userID
	err = c.Validate(request)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}
	order, err := h.ordersUsecase.InsertOrder(c.Request().Context(), request)
	if err != nil {
		statusCode := CreateOrderCustomErrorHTTPCode(err)
		response.Error = err.Error()
		return c.JSON(statusCode, response)
	}
	order.Result = true
	return c.JSON(http.StatusCreated, order)
}
