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

func TestHandler_CreateReview(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.PublicKey())

	reqJSON := `{"star": 5, "content": "Good!"}`
	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/reviews/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()

	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/reviews/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	err := jwtMiddleware(func(context echo.Context) error {
		return h.CreateReview(c)
	})(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestHandler_GetListOfReviewsByCourse(t *testing.T) {
	tearDown()
	setup()

	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/:id/reviews/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/reviews/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	err := h.GetListOfReviewsByCourse(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var rs reviewListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &rs)
		assert.NoError(t, err)
		assert.Equal(t, 1, rs.Count)
		assert.Equal(t, testReview.Content, rs.Reviews[0].Content)
	}
}

func TestHandler_GetListOfReviewsByUser(t *testing.T) {
	tearDown()
	setup()

	req := utils.NewTestRequest(echo.GET, "/api/v1/users/:id/reviews/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/users/:id/reviews/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testStudent.ID))
	err := h.GetListOfReviewsByUser(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var rs reviewListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &rs)
		assert.NoError(t, err)
		assert.Equal(t, 1, rs.Count)
		assert.Equal(t, testReview.Content, rs.Reviews[0].Content)
	}
}
