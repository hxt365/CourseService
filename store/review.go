package store

import (
	"CourseService/model"
	"github.com/jinzhu/gorm"
)

type ReviewStore struct {
	db *gorm.DB
}

func NewReviewStore(db *gorm.DB) *ReviewStore {
	return &ReviewStore{db: db}
}

func (rs *ReviewStore) Create(r *model.Review) error {
	return rs.db.Create(r).Error
}

func (rs *ReviewStore) ListByUser(userID uint, offset, limit int) ([]model.Review, int, error) {
	var (
		reviews []model.Review
		count   int
	)
	rs.db.Model(&reviews).Where(&model.Review{User: userID}).Count(&count)
	rs.db.Where(&model.Review{User: userID}).Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews)
	return reviews, count, nil
}

func (rs *ReviewStore) ListByCourse(courseID uint, offset, limit int) ([]model.Review, int, error) {
	var (
		course  model.Course
		reviews []model.Review
		count   int
	)
	if err := rs.db.Find(&course, courseID).Error; err != nil {
		return nil, 0, err
	}
	rs.db.Model(&course).Offset(offset).Limit(limit).Association("Reviews").Find(&reviews)
	count = rs.db.Model(&course).Association("Reviews").Count()
	return reviews, count, nil
}
