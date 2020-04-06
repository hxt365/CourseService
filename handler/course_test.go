package handler

import (
	"CourseService/router/middleware"
	"CourseService/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_GetCourse(t *testing.T) {
	tearDown()
	setup()
	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/:id/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	assert.NoError(t, h.GetCourse(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var c singleCourseResponse
		err := json.Unmarshal(rec.Body.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, testCourse.Name, c.Course.Name)
		assert.Equal(t, testCourse.MaxStudent, c.Course.MaxStudent)
	}
	cache, err := ca.GetCourseCache(testCourse.ID)
	if assert.NoError(t, err) {
		var c singleCourseResponse
		err := json.Unmarshal([]byte(cache), &c)
		assert.NoError(t, err)
		assert.Equal(t, testCourse.Name, c.Course.Name)
		assert.Equal(t, testCourse.MaxStudent, c.Course.MaxStudent)
	}
}

func TestHandler_GetListOfCourses(t *testing.T) {
	tearDown()
	setup()
	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	assert.NoError(t, h.GetListOfCourses(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var cs coursesListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cs)
		assert.NoError(t, err)
		assert.Equal(t, 2, cs.Count)
		assert.Equal(t, testCourse.Name, cs.Courses[0].Name)
	}

	cache, err := ca.GetCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT)
	if assert.NoError(t, err) {
		var cs coursesListResponse
		err := json.Unmarshal([]byte(cache), &cs)
		assert.NoError(t, err)
		assert.Equal(t, 2, cs.Count)
		assert.Equal(t, testCourse.Name, cs.Courses[0].Name)
	}
}

func TestHandler_GetListOfCoursesByMentor(t *testing.T) {
	tearDown()
	setup()
	req := utils.NewTestRequest(echo.GET, "/api/v1/mentors/:id/courses/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/mentors/:id/courses/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testMentorID))
	assert.NoError(t, h.GetListOfCoursesByMentor(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var cs coursesListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cs)
		assert.NoError(t, err)
		assert.Equal(t, 1, cs.Count)
		assert.Equal(t, testCourse.Name, cs.Courses[0].Name)
	}
	cache, err := ca.GetMentorsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testMentorID)
	if assert.NoError(t, err) {
		var cs coursesListResponse
		err := json.Unmarshal([]byte(cache), &cs)
		assert.NoError(t, err)
		assert.Equal(t, 1, cs.Count)
		assert.Equal(t, testCourse.Name, cs.Courses[0].Name)
	}
}

func TestHandler_GetListOfCoursesByStudent(t *testing.T) {
	tearDown()
	setup()
	req := utils.NewTestRequest(echo.GET, "/api/v1/students/:id/courses/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/students/:id/courses/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testStudent.ID))
	assert.NoError(t, h.GetListOfCoursesByStudent(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var cs coursesListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cs)
		assert.NoError(t, err)
		assert.Equal(t, 1, cs.Count)
		assert.Equal(t, testCourse.Name, cs.Courses[0].Name)
	}
	cache, err := ca.GetStudentsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testStudent.ID)
	if assert.NoError(t, err) {
		var cs coursesListResponse
		err := json.Unmarshal([]byte(cache), &cs)
		assert.NoError(t, err)
		assert.Equal(t, 1, cs.Count)
		assert.Equal(t, testCourse.Name, cs.Courses[0].Name)
	}
}

func TestHandler_CreateCourse(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"name": "Another course", "description": "Good course", "prerequisite": "nothing", "aim": "good programmer", "maxStudent": 5, "fee": 3}`
	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.CreateCourse(c)
	}))(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		var c singleCourseResponse
		err := json.Unmarshal(rec.Body.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, "Another course", c.Course.Name)
		assert.Equal(t, 5, int(c.Course.MaxStudent))
	}
	_, err = ca.GetCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT)
	assert.Error(t, err)
	_, err = ca.GetMentorsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testMentorID)
	assert.Error(t, err)
}

func TestHandler_CreateCourseByStudent(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"name": "Another course", "description": "Good course", "prerequisite": "nothing", "aim": "good programmer", "maxStudent": 5, "fee": 3}`
	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.CreateCourse(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_UpdateCourse(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"name": "Another course", "fee": 30}`
	req := utils.NewTestRequest(echo.PUT, "/api/v1/courses/:id/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.UpdateCourse(c)
	}))(c)
	assert.NoError(t, err)

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var c singleCourseResponse
		err := json.Unmarshal(rec.Body.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, "Another course", c.Course.Name)
		assert.Equal(t, 30, int(c.Course.Fee))
		assert.Equal(t, testCourse.Description, c.Course.Description)
	}
	_, err = ca.GetCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT)
	assert.Error(t, err)
	_, err = ca.GetCourseCache(testCourse.ID)
	assert.Error(t, err)
	_, err = ca.GetMentorsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testCourse.Mentor)
	assert.Error(t, err)
	_, err = ca.GetStudentsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testCourse.Students[0].ID)
	assert.Error(t, err)
}

func TestHandler_UpdateCourseOfAnother(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"name": "Another course", "fee": 30}`
	req := utils.NewTestRequest(echo.PUT, "/api/v1/courses/:id/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testAnotherMentorID, "mentor")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.UpdateCourse(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_UpdateCourseByStudent(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"name": "Another course", "fee": 30}`
	req := utils.NewTestRequest(echo.PUT, "/api/v1/courses/:id/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.UpdateCourse(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_DeleteCourse(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	req := utils.NewTestRequest(echo.DELETE, "/api/v1/courses/:id/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.DeleteCourse(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	_, err = ca.GetCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT)
	assert.Error(t, err)
	_, err = ca.GetCourseCache(testCourse.ID)
	assert.Error(t, err)
	_, err = ca.GetMentorsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testCourse.Mentor)
	assert.Error(t, err)
	_, err = ca.GetStudentsCoursesListCache(utils.DEFAULT_OFFSET, utils.DEFAULT_LIMIT, testCourse.Students[0].ID)
	assert.Error(t, err)
}

func TestHandler_DeleteCourseOfAnother(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	req := utils.NewTestRequest(echo.DELETE, "/api/v1/courses/:id/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testAnotherMentorID, "mentor")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.DeleteCourse(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_DeleteCourseByStudent(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	req := utils.NewTestRequest(echo.DELETE, "/api/v1/courses/:id/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(onlyMentorMiddleware(func(context echo.Context) error {
		return h.DeleteCourse(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_TakeCourse(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())

	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/take/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/take/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testAnotherCourse.ID))

	err := jwtMiddleware(func(context echo.Context) error {
		return h.TakeCourse(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	cache, err := ca.GetTakenCourseCache(testAnotherCourse.ID, testStudent.ID)
	assert.NoError(t, err)
	assert.Equal(t,true, cache)
}

func TestHandler_TakeCourseDuplicated(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())

	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/take/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/take/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(func(context echo.Context) error {
		return h.TakeCourse(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_TakeCourseByItsMentor(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())

	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/take/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/take/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(func(context echo.Context) error {
		return h.TakeCourse(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_DropCourse(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())

	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/drop/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/drop/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))

	err := jwtMiddleware(func(context echo.Context) error {
		return h.DropCourse(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	_, err = ca.GetTakenCourseCache(testCourse.ID, testStudent.ID)
	assert.Error(t, err)
}

func TestHandler_DropCourseNotTaken(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())

	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/drop/", nil)
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/drop/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testAnotherCourse.ID))

	err := jwtMiddleware(func(context echo.Context) error {
		return h.DropCourse(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandle_GetCourseFromCache(t *testing.T) {
	tearDown()
	setup()

	// Get first time, query from db
	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/:id/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	assert.NoError(t, h.GetCourse(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	// Get second time, query from cache
	req2 := utils.NewTestRequest(echo.GET, "/api/v1/courses/:id/", nil)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2.Request, rec2)
	c2.SetPath("/api/v1/courses/:id/")
	c2.SetParamNames("id")
	c2.SetParamValues(fmt.Sprint(testCourse.ID))
	assert.NoError(t, h.GetCourse(c2))
	if assert.Equal(t, http.StatusOK, rec2.Code) {
		fmt.Println(rec2.Body.String())
		var c singleCourseResponse
		err := json.Unmarshal(rec2.Body.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, testCourse.Name, c.Course.Name)
		assert.Equal(t, testCourse.MaxStudent, c.Course.MaxStudent)
	}
}