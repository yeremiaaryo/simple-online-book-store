package users

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/users"
)

//go:generate mockgen -package=users -source=users_handler.go -destination=users_handler_mock_test.go
type usersUsecase interface {
	CreateUser(ctx context.Context, req users.CreateUserRequest) (*users.Model, error)
	Login(ctx context.Context, req users.LoginRequest) (string, error)
}
type Handler struct {
	usersUsecase usersUsecase
}

func New(usersUsecase usersUsecase) *Handler {
	return &Handler{usersUsecase: usersUsecase}
}

func (h *Handler) CreateUser(c echo.Context) error {
	response := users.UserResponse{}
	var request users.CreateUserRequest
	err := c.Bind(&request)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	err = c.Validate(request)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}
	user, err := h.usersUsecase.CreateUser(c.Request().Context(), request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.EqualFold(err.Error(), "email already exists") {
			statusCode = http.StatusFound
		}
		response.Error = err.Error()
		return c.JSON(statusCode, response)
	}
	response.Result = true
	response.User = user
	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(c echo.Context) error {
	response := users.LoginResponse{}
	var request users.LoginRequest
	err := c.Bind(&request)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	err = c.Validate(request)
	if err != nil {
		response.Error = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	token, err := h.usersUsecase.Login(c.Request().Context(), request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid email") {
			statusCode = http.StatusUnauthorized
		}
		response.Error = err.Error()
		return c.JSON(statusCode, response)
	}
	response.Result = true
	response.Token = token
	return c.JSON(http.StatusOK, response)
}
