package util

import (
	"errors"
	"github.com/labstack/echo/v4"
)

func GetUserID(c echo.Context) (int64, error) {
	userID := c.Get("userID")
	if userID == nil {
		return 0, errors.New("userID not found")
	}
	return userID.(int64), nil
}
