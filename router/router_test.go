package router

import (
	"CourseService/router/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddAuthHeader(t *testing.T) {
	addAuthHeaderMiddleware := middleware.AddAuthHeader()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:     "access-token",
		Value:    "sometoken",
		Expires:  time.Now().Add(time.Minute * 5),
		HttpOnly: true,
	})
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	err := addAuthHeaderMiddleware(func(context echo.Context) error {
		return nil
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, "JWT sometoken", req.Header.Get("Authorization"))
}

func TestAddAuthHeaderAlreadyHad(t *testing.T) {
	addAuthHeaderMiddleware := middleware.AddAuthHeader()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set("Authorization", "JWT oldtoken")
	req.AddCookie(&http.Cookie{
		Name:     "access-token",
		Value:    "sometoken",
		Expires:  time.Now().Add(time.Minute * 5),
		HttpOnly: true,
	})
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	err := addAuthHeaderMiddleware(func(context echo.Context) error {
		return nil
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, "JWT oldtoken", req.Header.Get("Authorization"))
}

func TestAddAuthHeaderNoCookieNoHeader(t *testing.T) {
	addAuthHeaderMiddleware := middleware.AddAuthHeader()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	err := addAuthHeaderMiddleware(func(context echo.Context) error {
		return nil
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, "", req.Header.Get("Authorization"))
}
