package handler

import "CourseService/model"

type courseResponse struct {
	ID           uint              `json:"id"`
	Mentor       uint              `json:"tutor"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Prerequisite string            `json:"prerequisite"`
	Aim          string            `json:"aim"`
	MaxStudent   uint              `json:"maxStudent"`
	Fee          uint              `json:"fee"`
	Rating       uint              `json:"rating"`
	Reviews      []*reviewResponse `json:"reviews"`
	CreatedAt    string            `json:"createdAt"`
	UpdatedAt    string            `json:"updatedAt"`
}

type singleCourseResponse struct {
	Course *courseResponse `json:"course"`
}

type courseListResponse struct {
	Courses []*courseResponse `json:"courses"`
	Count   int               `json:"count"`
}

func newCourseResponse(course *model.Course) *singleCourseResponse {
	c := &courseResponse{
		ID:           course.ID,
		Mentor:       course.Mentor,
		Name:         course.Name,
		Description:  course.Description,
		Prerequisite: course.Prerequisite,
		Aim:          course.Aim,
		MaxStudent:   course.MaxStudent,
		Fee:          course.Fee,
		Rating:       course.Rating,
		CreatedAt:    course.CreatedAt.String(),
		UpdatedAt:    course.UpdatedAt.String(),
	}
	for _, review := range course.Reviews {
		c.Reviews = append(c.Reviews, &reviewResponse{
			User:      review.User,
			Star:      review.Star,
			Content:   review.Content,
			CreatedAt: review.CreatedAt.String(),
		})
	}
	return &singleCourseResponse{c}
}

func newCourseListResponse(courses []model.Course, count int) *courseListResponse {
	cs := &courseListResponse{
		Count: count,
	}
	for _, c := range courses {
		cs.Courses = append(cs.Courses, &courseResponse{
			ID:           c.ID,
			Mentor:       c.Mentor,
			Name:         c.Name,
			Description:  c.Description,
			Prerequisite: c.Prerequisite,
			Aim:          c.Aim,
			MaxStudent:   c.MaxStudent,
			Fee:          c.Fee,
			Rating:       c.Rating,
			CreatedAt:    c.CreatedAt.String(),
			UpdatedAt:    c.UpdatedAt.String(),
		})
	}
	return cs
}

type reviewResponse struct {
	User      uint   `json:"user"`
	Star      uint   `json:"star"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type reviewListResponse struct {
	Reviews []*reviewResponse `json:"reviews"`
	Count   int               `json:"count"`
}

func newReviewListResponse(reviews []model.Review, count int) *reviewListResponse {
	rs := &reviewListResponse{
		Count: count,
	}
	for _, r := range reviews {
		rs.Reviews = append(rs.Reviews, &reviewResponse{
			User:      r.User,
			Star:      r.Star,
			Content:   r.Content,
			CreatedAt: r.CreatedAt.String(),
		})
	}
	return rs
}

type notificationResponse struct {
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type notificationListResponse struct {
	Notifications []*notificationResponse `json:"notifications"`
	Count         int                     `json:"count"`
}

func newNotificationListResponse(notifications []model.Notification, count int) *notificationListResponse {
	ns := &notificationListResponse{
		Count: count,
	}
	for _, n := range notifications {
		ns.Notifications = append(ns.Notifications, &notificationResponse{
			Content:   n.Content,
			CreatedAt: n.CreatedAt.String(),
		})
	}
	return ns
}
