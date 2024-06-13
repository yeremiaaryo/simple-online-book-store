package middleware

import (
	"net/http"
	"strings"

	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/pkg/jwt"

	"github.com/yeremiaaryo/gotu-assignment/internal/response"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
		c.Set("userID", userID)
		return next(c)
	}
}
