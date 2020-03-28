package handler

import (
	"CourseService/model"
	"CourseService/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// Only tutor
func (h *Handler) CreateCourse(c echo.Context) error {
	var course model.Course
	req := &courseCreateRequest{}
	if err := req.bind(c, &course); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	course.Mentor = c.Get("user").(uint)
	if err := h.courseStore.Create(&course); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newCourseResponse(&course))
}

func (h *Handler) GetCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	course, err := h.courseStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newCourseResponse(course))
}

// Only tutor
func (h *Handler) UpdateCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	course, err := h.courseStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if course.Mentor != c.Get("user") {
		return c.JSON(http.StatusForbidden, utils.AccessForbiden())
	}
	req := &courseUpdateRequest{}
	req.populate(course)
	if err := req.bind(c, course); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.courseStore.Update(course); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newCourseResponse(course))
}

// Only tutor
func (h *Handler) DeleteCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	course, err := h.courseStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	userID := c.Get("user").(uint)
	if userID != course.Mentor {
		return c.JSON(http.StatusForbidden, utils.AccessForbiden())
	}
	if err := h.courseStore.Delete(course); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (h *Handler) GetListOfCourses(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	courses, count, err := h.courseStore.List(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newCourseListResponse(courses, count))
}

func (h *Handler) GetListOfCoursesByMentor(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	courses, count, err := h.courseStore.ListByTutor(uint(id), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newCourseListResponse(courses, count))
}

func (h *Handler) GetListOfCoursesByStudent(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	courses, count, err := h.courseStore.ListByStudent(uint(id), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newCourseListResponse(courses, count))
}

// Authenticated
func (h *Handler) TakeCourse(c echo.Context) error {
	userID := c.Get("user").(uint)
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	course, err := h.courseStore.GetByID(uint(courseID))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if course.Mentor == userID {
		return c.JSON(http.StatusBadRequest, utils.BadRequest())
	}
	if h.courseStore.IfStudentTookCourse(userID, uint(courseID)) {
		return c.JSON(http.StatusBadRequest, utils.BadRequest())
	}
	if err := h.courseStore.CourseTakenByStudent(course, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

// Authenticated
func (h *Handler) DropCourse(c echo.Context) error {
	userID := c.Get("user").(uint)
	courseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	course, err := h.courseStore.GetByID(uint(courseID))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if !h.courseStore.IfStudentTookCourse(userID, uint(courseID)) {
		return c.JSON(http.StatusBadRequest, utils.BadRequest())
	}
	if err := h.courseStore.DropCourseByStudent(course, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}
