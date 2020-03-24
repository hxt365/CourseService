package utils

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func GetOffsetLimit(c echo.Context) (int, int) {
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	return offset, limit
}
