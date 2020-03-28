package utils

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
)

type TestRequest struct {
	Request *http.Request
}

func NewTestRequest(method, target string, body io.Reader) *TestRequest {
	return &TestRequest{Request: httptest.NewRequest(method, target, body)}
}

func (req *TestRequest) SetAuthHeader(id uint, role string) {
	req.Request.Header.Set(echo.HeaderAuthorization, authHeader(generateJWTToken(id, role)))
}

func (req *TestRequest) SetJSONHeader() {
	req.Request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}

func authHeader(token string) string {
	return "JWT " + token
}
