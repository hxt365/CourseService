package middleware

import (
	"CourseService/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func OnlyTutor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role")
			if role == "tutor" {
				return next(c)
			}
			return c.JSON(http.StatusForbidden, utils.AccessForbiden())
		}
	}
}
