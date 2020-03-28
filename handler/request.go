package handler

import (
	"CourseService/model"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type courseCreateRequest struct {
	Name         string `json:"name" validate:"required,max=100"`
	Description  string `json:"description" validate:"required,max=300"`
	Prerequisite string `json:"prerequisite" validate:"required,max=300"`
	Aim          string `json:"aim" validate:"required,max=300"`
	MaxStudent   uint   `json:"maxStudent" validate:"required,min=1,max=10"`
	Fee          uint   `json:"fee" validate:"required,min=0"`
}

func (r *courseCreateRequest) bind(c echo.Context, course *model.Course) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	_ = copier.Copy(course, r)
	return nil
}

type courseUpdateRequest struct {
	Name         string `json:"name" validate:"max=100"`
	Description  string `json:"description" validate:"max=300"`
	Prerequisite string `json:"prerequisite" validate:"max=300"`
	Aim          string `json:"aim" validate:"max=300"`
	MaxStudent   uint   `json:"maxStudent" validate:"min=1,max=10"`
	Fee          uint   `json:"fee" validate:"min=0"`
}

func (r *courseUpdateRequest) populate(course *model.Course) {
	r.Name = course.Name
	r.Description = course.Description
	r.Prerequisite = course.Prerequisite
	r.Aim = course.Aim
	r.MaxStudent = course.MaxStudent
	r.Fee = course.Fee
}

func (r *courseUpdateRequest) bind(c echo.Context, course *model.Course) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	_ = copier.Copy(course, r)
	return nil
}

type reviewCreateRequest struct {
	Star    uint   `json:"star" validate:"required,min=1,max=5"`
	Content string `json:"content" validate:"required,max=300"`
}

func (r *reviewCreateRequest) bind(c echo.Context, review *model.Review) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	_ = copier.Copy(review, r)
	return nil
}

type notificationCreateRequest struct {
	Content string `json:"content" validate:"required,max=300"`
}

func (r *notificationCreateRequest) bind(c echo.Context, notification *model.Notification) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	notification.Content = r.Content
	return nil
}
