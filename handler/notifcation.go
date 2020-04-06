package handler

import "C"
import (
	"CourseService/model"
	"CourseService/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// only tutor
func (h *Handler) CreateNotification(c echo.Context) error {
	courseID, err := strconv.Atoi(c.Param("id"))
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

	if course.Mentor != c.Get("user").(uint) {
		return c.JSON(http.StatusForbidden, utils.AccessForbiden())
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
	h.cacheStore.DeleteCoursesNotificationsListCache(uint(courseID))
	return c.JSON(http.StatusCreated, map[string]interface{}{"result": "ok"})
}

// Authenticated
func (h *Handler) GetListOfNotifications(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	userID := c.Get("user").(uint)

	if !h.courseStore.IfStudentTookCourse(userID, uint(courseID)) {
		return c.JSON(http.StatusForbidden, utils.AccessForbiden())
	}

	cached, err := h.cacheStore.GetCoursesNotificationsListCache(offset, limit, uint(courseID))
	if err == nil {
		var res notificationsListResponse
		err := json.Unmarshal([]byte(cached), &res)
		if err == nil {
			return c.JSON(http.StatusOK, res)
		}
	}

	notifications, count, err := h.notificationStore.ListByCourse(uint(courseID), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	res := newNotificationListResponse(notifications, count)
	resBytes, err := json.Marshal(res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	h.cacheStore.SetCoursesNotificationsListCache(offset, limit, uint(courseID), string(resBytes))
	return c.JSON(http.StatusOK, res)
}
