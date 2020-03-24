package handler

import (
	"CourseService/model"
	"CourseService/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// Authenticated
func (h *Handler) CreateReview(c echo.Context) error {
	userID := c.Get("user").(uint)
	courseID, err := strconv.Atoi(c.Param("course"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	course, err := h.courseStore.GetByID(uint(courseID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	var review model.Review
	req := &reviewCreateRequest{}
	if err := req.bind(c, &review); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	review.User = userID
	review.CourseID = uint(courseID)
	if err := h.reviewStore.Create(&review); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"result": "ok"})
}

func (h *Handler) GetListOfReviewsByCourse(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	courseID, err := strconv.Atoi(c.Param("course"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	reviews, count, err := h.reviewStore.ListByCourse(uint(courseID), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newReviewListResponse(reviews, count))
}

func (h *Handler) GetListOfReviewsByUser(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	userID, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	reviews, count, err := h.reviewStore.ListByUser(uint(userID), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newReviewListResponse(reviews, count))
}
