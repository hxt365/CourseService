package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func AddAuthHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if len(authHeader) == 0 {
				authCookie, err := c.Request().Cookie("access-token")
				if err == nil {
					c.Request().Header.Set("Authorization", fmt.Sprint("JWT ", authCookie.Value))
				}
			}
			return next(c)
		}
	}
}
