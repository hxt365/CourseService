package store

import (
	"CourseService/model"
	"github.com/jinzhu/gorm"
)

type CourseStore struct {
	db *gorm.DB
}

func NewCourseStore(db *gorm.DB) *CourseStore {
	return &CourseStore{db: db}
}

func (cs *CourseStore) Create(c *model.Course) error {
	return cs.db.Create(c).Error
}

func (cs *CourseStore) GetByID(id uint) (*model.Course, error) {
	var c model.Course
	err := cs.db.Preload("Reviews", func(db *gorm.DB) *gorm.DB {
		return db.Offset(0).Limit(10).Order("created_at DESC")
	}).First(&c, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (cs *CourseStore) Update(c *model.Course) error {
	return cs.db.Model(c).Select("name, description, prerequisite, aim, maxStudent, fee").Update(c).Error
}

func (cs *CourseStore) Delete(c *model.Course) error {
	return cs.db.Delete(c).Error
}

func (cs *CourseStore) List(offset, limit int) ([]model.Course, int, error) {
	var (
		courses []model.Course
		count   int
	)
	cs.db.Model(&courses).Count(&count)
	cs.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&courses)
	return courses, count, nil
}

func (cs *CourseStore) ListByTutor(tutorID uint, offset, limit int) ([]model.Course, int, error) {
	var (
		courses []model.Course
		count   int
	)
	cs.db.Where(&model.Course{Mentor: tutorID}).Offset(offset).Limit(limit).Order("created_at DESC").Find(&courses)
	cs.db.Model(&model.Course{}).Where(&model.Course{Mentor: tutorID}).Count(&count)
	return courses, count, nil
}

func (cs *CourseStore) ListByStudent(studentID uint, offset, limit int) ([]model.Course, int, error) {
	var (
		student model.Student
		courses []model.Course
		count   int
	)
	if err := cs.db.First(&student, studentID).Error; err != nil {
		return nil, 0, err
	}
	cs.db.Model(&student).Offset(offset).Limit(limit).Order("created_at DESC").Association("Courses").
		Find(&courses)
	count = cs.db.Model(&student).Association("Courses").Count()
	return courses, count, nil
}

func (cs *CourseStore) IfStudentTookCourse(studentID, courseID uint) bool {
	var count int
	cs.db.Table("student_course").Select("1").
		Where("student_id = ? AND course_id = ?", studentID, courseID).Count(&count)
	return count > 0
}

func (cs *CourseStore) CourseTakenByStudent(course *model.Course, userID uint) error {
	var student model.Student
	if err := cs.db.FirstOrCreate(&student, model.Student{ID: userID,}).Error; err != nil {
		return err
	}
	if err := cs.db.Model(course).Association("Students").Append(&student).Error; err != nil {
		return err
	}
	return nil
}

func (cs *CourseStore) DropCourseByStudent(course *model.Course, userID uint) error {
	var student model.Student
	if err := cs.db.First(&student, userID).Error; err != nil {
		return err
	}
	if err := cs.db.Model(course).Association("Students").Delete(&student).Error; err != nil {
		return err
	}
	return nil
}
