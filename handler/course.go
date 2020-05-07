package handler

import (
	"CourseService/model"
	"CourseService/utils"
	"encoding/json"
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
	userID := c.Get("user").(uint)
	course.Mentor = userID
	if err := h.courseStore.Create(&course); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	h.cacheStore.DeleteCoursesListCache()
	h.cacheStore.DeleteMentorsCoursesListCache(userID)
	return c.JSON(http.StatusCreated, newCourseResponse(&course))
}

func (h *Handler) GetCourse(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	cached, err := h.cacheStore.GetCourseCache(uint(id))
	if err == nil {
		var res singleCourseResponse
		err := json.Unmarshal([]byte(cached), &res)
		if err == nil {
			return c.JSON(http.StatusOK, res)
		}
	}

	course, err := h.courseStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if course == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	res := newCourseResponse(course)
	resBytes, err := json.Marshal(res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	h.cacheStore.SetCourseCache(uint(id), string(resBytes))
	return c.JSON(http.StatusOK, res)
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

	h.cacheStore.DeleteCoursesListCache()
	h.cacheStore.DeleteCourseCache(uint(id))
	h.cacheStore.DeleteMentorsCoursesListCache(course.Mentor)
	h.cacheStore.DeleteAllStudentsCoursesListCacheByCourse(course)
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

	h.cacheStore.DeleteCoursesListCache()
	h.cacheStore.DeleteCourseCache(uint(id))
	h.cacheStore.DeleteMentorsCoursesListCache(course.Mentor)
	h.cacheStore.DeleteAllStudentsCoursesListCacheByCourse(course)
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (h *Handler) GetListOfCourses(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)

	cached, err := h.cacheStore.GetCoursesListCache(offset, limit)
	if err == nil {
		var res coursesListResponse
		err := json.Unmarshal([]byte(cached), &res)
		if err == nil {
			return c.JSON(http.StatusOK, res)
		}
	}

	courses, count, err := h.courseStore.List(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	res := newCourseListResponse(courses, count)
	resBytes, err := json.Marshal(res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	h.cacheStore.SetCoursesListCache(offset, limit, string(resBytes))
	return c.JSON(http.StatusOK, newCourseListResponse(courses, count))
}

func (h *Handler) GetListOfCoursesByMentor(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	cached, err := h.cacheStore.GetMentorsCoursesListCache(offset, limit, uint(id))
	if err == nil {
		var res coursesListResponse
		err := json.Unmarshal([]byte(cached), &res)
		if err == nil {
			return c.JSON(http.StatusOK, res)
		}
	}

	courses, count, err := h.courseStore.ListByTutor(uint(id), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	res := newCourseListResponse(courses, count)
	resBytes, err := json.Marshal(res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	h.cacheStore.SetMentorsCoursesListCache(offset, limit, uint(id), string(resBytes))
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetListOfCoursesByStudent(c echo.Context) error {
	offset, limit := utils.GetOffsetLimit(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	cached, err := h.cacheStore.GetStudentsCoursesListCache(offset, limit, uint(id))
	if err == nil {
		var res coursesListResponse
		err := json.Unmarshal([]byte(cached), &res)
		if err == nil {
			return c.JSON(http.StatusOK, res)
		}
	}

	courses, count, err := h.courseStore.ListByStudent(uint(id), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	res := newCourseListResponse(courses, count)
	resBytes, err := json.Marshal(res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	h.cacheStore.SetStudentsCoursesListCache(offset, limit, uint(id), string(resBytes))

	return c.JSON(http.StatusOK, res)
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

	cached, err := h.cacheStore.GetTakenCourseCache(uint(courseID), userID)
	if err == nil && cached {
		return c.JSON(http.StatusBadRequest, utils.BadRequest())
	}
	if err != nil && h.courseStore.IfStudentTookCourse(userID, uint(courseID)) {
		return c.JSON(http.StatusBadRequest, utils.BadRequest())
	}

	if err := h.courseStore.CourseTakenByStudent(course, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	h.cacheStore.SetTakenCourseCache(uint(courseID), userID, true)
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

	h.cacheStore.DeleteTakenCourseCache(uint(courseID), userID)
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}
