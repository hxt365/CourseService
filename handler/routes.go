package handler

import (
	"CourseService/router/middleware"
	"github.com/labstack/echo/v4"
	"os"
)

func (h *Handler) Register(r *echo.Group) {
	jwtMiddleware := middleware.JWT(os.Getenv("JWT_PUBLIC"))
	onlyTutorMiddleware := middleware.OnlyTutor()

	courses := r.Group("courses/")
	courses.GET("", h.GetListOfCourses)
	courses.GET(":id/", h.GetCourse)
	courses.GET(":id/reviews/", h.GetListOfReviewsByCourse)

	courses.POST(":id/take/", h.TakeCourse, jwtMiddleware)
	courses.POST(":id/drop/", h.DropCourse, jwtMiddleware)
	courses.POST(":id/reviews/", h.CreateReview, jwtMiddleware)
	courses.POST(":id/notifications/", h.CreateNotification, jwtMiddleware)
	courses.GET(":id/notifications/", h.GetListOfNotifications, jwtMiddleware)

	courses.POST("", h.CreateCourse, jwtMiddleware, onlyTutorMiddleware)
	courses.PUT("/:id/", h.UpdateCourse, jwtMiddleware, onlyTutorMiddleware)
	courses.PATCH("/:id/", h.UpdateCourse, jwtMiddleware, onlyTutorMiddleware)
	courses.DELETE(":id/", h.DeleteCourse, jwtMiddleware, onlyTutorMiddleware)

	users := r.Group("users/")
	users.GET(":id/reviews/", h.GetListOfReviewsByUser)

	tutors := r.Group("tutors/")
	tutors.GET(":id/courses/", h.GetListOfCoursesByTutor)

	students := r.Group("students/")
	students.GET(":id/courses/", h.GetListOfCoursesByStudent)
}
