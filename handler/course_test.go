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
}

func TestHandler_GetListOfCourses(t *testing.T) {
	tearDown()
	setup()
	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	assert.NoError(t, h.GetListOfCourses(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var cs courseListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cs)
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
		var cs courseListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cs)
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
		var cs courseListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cs)
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
