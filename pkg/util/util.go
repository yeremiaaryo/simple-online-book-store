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

// GetLimitAndOffset calculates and returns the limit and offset based on the provided pageIndex and pageSize.
func GetLimitAndOffset(pageIndex, pageSize int) (int, int) {
	// Sanitize request
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Convert to limit and offset
	limit := pageSize
	offset := (pageIndex - 1) * limit

	return limit, offset
}
