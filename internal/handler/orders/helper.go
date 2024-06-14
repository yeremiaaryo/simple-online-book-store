package orders

import (
	"net/http"
	"strings"
)

func CreateOrderCustomErrorHTTPCode(err error) int {
	if strings.Contains(err.Error(), "book with id") || strings.Contains(err.Error(), "total amount is different") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
