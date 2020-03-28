package handler

import (
	"CourseService/router/middleware"
	"CourseService/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(r *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.PublicKey())
	onlyMentorMiddleware := middleware.OnlyMentor()

	courses := r.Group("courses/")
	courses.GET("", h.GetListOfCourses, jwtMiddleware)
	courses.GET(":id/", h.GetCourse)
	courses.GET(":id/reviews/", h.GetListOfReviewsByCourse)

	courses.POST(":id/take/", h.TakeCourse, jwtMiddleware)
	courses.POST(":id/drop/", h.DropCourse, jwtMiddleware)
	courses.POST(":id/reviews/", h.CreateReview, jwtMiddleware)
	courses.POST(":id/notifications/", h.CreateNotification, jwtMiddleware)
	courses.GET(":id/notifications/", h.GetListOfNotifications, jwtMiddleware)

	courses.POST("", h.CreateCourse, jwtMiddleware, onlyMentorMiddleware)
	courses.PUT("/:id/", h.UpdateCourse, jwtMiddleware, onlyMentorMiddleware)
	courses.PATCH("/:id/", h.UpdateCourse, jwtMiddleware, onlyMentorMiddleware)
	courses.DELETE(":id/", h.DeleteCourse, jwtMiddleware, onlyMentorMiddleware)

	users := r.Group("users/")
	users.GET(":id/reviews/", h.GetListOfReviewsByUser)

	mentors := r.Group("mentors/")
	mentors.GET(":id/courses/", h.GetListOfCoursesByMentor)

	students := r.Group("students/")
	students.GET(":id/courses/", h.GetListOfCoursesByStudent)
}
