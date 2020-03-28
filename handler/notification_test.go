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

func TestHandler_CreateNotification(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onTutorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"content": "Good!"}`
	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/notifications/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("api/v1/courses/:id/notifications/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	err := jwtMiddleware(onTutorMiddleware(func(context echo.Context) error {
		return h.CreateNotification(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestHandler_CreateNotificationForAnotherCourse(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onTutorMiddleware := middleware.OnlyMentor()

	reqJSON := `{"content": "Good!"}`
	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/notifications/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("api/v1/courses/:id/notifications/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testAnotherCourse.ID))
	err := jwtMiddleware(onTutorMiddleware(func(context echo.Context) error {
		return h.CreateNotification(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_CreateNotificationByStudent(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onTutorMiddleware := middleware.OnlyMentor()
	reqJSON := `{"content": "Good!"}`
	req := utils.NewTestRequest(echo.POST, "/api/v1/courses/:id/notifications/", strings.NewReader(reqJSON))
	req.SetJSONHeader()
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("api/v1/courses/:id/notifications/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	err := jwtMiddleware(onTutorMiddleware(func(context echo.Context) error {
		return h.CreateNotification(c)
	}))(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetListOfNotifications(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/:id/notifications/", nil)
	req.SetAuthHeader(testStudent.ID, "student")
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/notifications/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetListOfNotifications(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var ns notificationListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &ns)
		assert.NoError(t, err)
		assert.Equal(t, 1, ns.Count)
		assert.Equal(t, testNotification.Content, ns.Notifications[0].Content)
	}
}

func TestHandler_GetListOfNotificationsNotByStudent(t *testing.T) {
	tearDown()
	setup()

	jwtMiddleware := middleware.JWT(utils.PublicKey())
	req := utils.NewTestRequest(echo.GET, "/api/v1/courses/:id/notifications/", nil)
	req.SetAuthHeader(testMentorID, "mentor")
	rec := httptest.NewRecorder()
	c := e.NewContext(req.Request, rec)
	c.SetPath("/api/v1/courses/:id/notifications/")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(testCourse.ID))
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetListOfNotifications(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
