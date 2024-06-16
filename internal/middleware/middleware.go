package middleware

import (
	"fmt"
	"github.com/yeremiaaryo/gotu-assignment/internal/constant"
	"net/http"
	"strings"

	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/pkg/jwt"

	"github.com/yeremiaaryo/gotu-assignment/internal/response"

	"github.com/labstack/echo/v4"
)

type redis interface {
	Get(key string, field ...interface{}) (string, error)
}

type Handler struct {
	redis redis
}

func New(redis redis) *Handler {
	return &Handler{redis: redis}
}

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(c echo.Context) error {
		var header = c.Request().Header.Get("Authorization") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			// Token is missing
			return c.JSON(http.StatusForbidden, response.BaseResponse{
				Result: false,
				Error:  "no token provided",
			})
		}
		tokenString := header[len("Bearer "):]
		userID, err := jwt.VerifyToken(tokenString, secretKey)
		if err != nil {
			// Token is invalid
			return c.JSON(http.StatusForbidden, response.BaseResponse{
				Result: false,
				Error:  "invalid token",
			})
		}
		latestToken, err := h.redis.Get(fmt.Sprintf(constant.RedisKeyToken, userID))
		if err != nil {
			return c.JSON(http.StatusForbidden, response.BaseResponse{
				Result: false,
				Error:  "invalid to retrieve token",
			})
		}
		if latestToken != tokenString {
			return c.JSON(http.StatusForbidden, response.BaseResponse{
				Result: false,
				Error:  "token is invalid, please re-login",
			})
		}
		c.Set("userID", userID)
		return next(c)
	}
}
