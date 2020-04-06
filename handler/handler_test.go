package handler

import (
	"CourseService/db"
	"CourseService/model"
	"CourseService/router"
	"CourseService/store"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
	"testing"
)

var (
	d                   *gorm.DB
	rd                  *db.Redis
	ca                  *store.CacheStore
	cs                  *store.CourseStore
	ns                  *store.NotificationStore
	rs                  *store.ReviewStore
	h                   *Handler
	e                   *echo.Echo
	testMentorID        uint
	testAnotherMentorID uint
	testCourse          model.Course
	testStudent         model.Student
	testReview          model.Review
	testNotification    model.Notification
	testAnotherCourse   model.Course
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	d = db.NewTestDB()
	db.AutoMigrate(d)

	rd = db.NewTestRedis()

	cs = store.NewCourseStore(d)
	ns = store.NewNotificationStore(d)
	rs = store.NewReviewStore(d)
	ca = store.NewCacheStore(rd)

	h = NewHandler(cs, ns, rs, ca)

	e = router.New()
	_ = loadFixtures()
}

func tearDown() {
	if err := db.DropTestDB(d); err != nil {
		log.Fatal(err)
	}
	if err := d.Exec(fmt.Sprint("CREATE DATABASE ", os.Getenv("TEST_DATABASE_NAME"))).Error; err != nil {
		log.Fatal(err)
	}
	_ = d.Close()
	defer rd.Client.(*miniredis.Miniredis).Close()
}

func loadFixtures() error {
	testMentorID = 1
	testAnotherMentorID = 2

	testCourse = model.Course{
		Mentor:       testMentorID,
		Name:         "Giao trinh tan gai",
		Description:  "Tan phat do luon",
		Prerequisite: "Dep trai + Nhieu tien",
		Aim:          "Tan bat ki ai",
		MaxStudent:   5,
		Fee:          10,
	}
	_ = cs.Create(&testCourse)

	testAnotherCourse = model.Course{
		Mentor:       testAnotherMentorID,
		Name:         "Giao trinh tan gai",
		Description:  "Tan phat do luon",
		Prerequisite: "Dep trai + Nhieu tien",
		Aim:          "Tan bat ki ai",
		MaxStudent:   5,
		Fee:          10,
	}
	_ = cs.Create(&testAnotherCourse)

	testStudent = model.Student{ID: 3,}
	_ = cs.CourseTakenByStudent(&testCourse, testStudent.ID)

	testReview = model.Review{
		User:     testStudent.ID,
		Star:     5,
		Content:  "Rat tot, giao vien dep trai va nhieu tien",
		CourseID: testCourse.ID,
	}
	_ = rs.Create(&testReview)

	testNotification = model.Notification{
		Content:  "Hay tro nen giau co!",
		CourseID: testCourse.ID,
	}
	_ = ns.Create(&testNotification)
	return nil
}
