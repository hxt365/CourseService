package utils

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["Body"] = v.Message
	default:
		e.Errors["Body"] = err.Error()
	}
	return e
}

func NewValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		e.Errors[v.Field()] = fmt.Sprintf("%v", v.Tag())
	}
	return e
}

func AccessForbiden() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["Body"] = "access forbiden"
	return e
}

func NotFound() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["Body"] = "resource not found"
	return e
}
