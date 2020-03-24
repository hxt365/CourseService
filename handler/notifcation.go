package handler

import "C"
import (
	"CourseService/model"
	"CourseService/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// only tutor
func (h *Handler) CreateNotification(c echo.Context) error {
	courseID, err := strconv.Atoi(c.Param("course"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	course, err := h.courseStore.GetByID(uint(courseID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	var notification model.Notification
	req := &notificationCreateRequest{}
	if err := req.bind(c, &notification); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	notification.CourseID = uint(courseID)
	if err := h.notificationStore.Create(&notification); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"result": "ok"})
}

// Authenticated
func (h *Handler) GetListOfNotifications(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	courseID, err := strconv.Atoi(c.Param("course"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	userID := c.Get("user").(uint)
	if !h.courseStore.IfStudentHasCourse(userID, uint(courseID)) {
		return c.JSON(http.StatusForbidden, utils.AccessForbiden())
	}
	notifications, count, err := h.notificationStore.ListByCourse(uint(courseID), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newNotificationListResponse(notifications, count))
}