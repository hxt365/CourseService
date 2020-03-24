package main

import (
	"CourseService/db"
	"CourseService/handler"
	"CourseService/router"
	"CourseService/store"
)

func main() {
	r := router.New()
	v1 := r.Group("/api/v1/")

	d := db.New()
	db.AutoMigrate(d)

	cs := store.NewCourseStore(d)
	ns := store.NewNotificationStore(d)
	rs := store.NewReviewStore(d)

	h := handler.NewHandler(cs, ns, rs)

	h.Register(v1)

	r.Logger.Fatal(r.Start("0.0.0.0:8080"))
}
