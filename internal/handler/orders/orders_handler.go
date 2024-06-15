package orders

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
	"github.com/yeremiaaryo/gotu-assignment/pkg/util"
	"net/http"
	"strconv"
)

//go:generate mockgen -package=orders -source=orders_handler.go -destination=orders_handler_mock_test.go
type ordersUsecase interface {
	InsertOrder(ctx context.Context, order orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
	GetOrdersByUserID(ctx context.Context, userID int64, pageIndex, pageSize int) ([]orders.History, error)
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

func (h *Handler) GetOrderHistory(c echo.Context) error {
	var response orders.OrderHistoryResponse

	userID, err := util.GetUserID(c)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	pageIndex, err := strconv.Atoi(c.QueryParam("page_index"))
	if err != nil {
		pageIndex = 1 // default page index is 1 if error
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		pageSize = 10 // default page size is 10 if error
	}

	orderHistory, err := h.ordersUsecase.GetOrdersByUserID(c.Request().Context(), userID, pageIndex, pageSize)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Histories = orderHistory
	response.Result = true
	return c.JSON(http.StatusOK, response)
}
