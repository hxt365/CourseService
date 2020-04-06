package utils

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func GetOffsetLimit(c echo.Context) (int, int) {
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = DEFAULT_OFFSET
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = DEFAULT_LIMIT
	}
	return offset, limit
}
