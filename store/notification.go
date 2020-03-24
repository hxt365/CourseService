package store

import (
	"CourseService/model"
	"github.com/jinzhu/gorm"
)

type NotificationStore struct {
	db *gorm.DB
}

func NewNotificationStore(db *gorm.DB) *NotificationStore {
	return &NotificationStore{db: db}
}

func (ns *NotificationStore) Create(notification *model.Notification) error {
	return ns.db.Create(notification).Error
}

func (ns *NotificationStore) ListByCourse(courseID uint, offset, limit int) ([]model.Notification, int, error) {
	var (
		course        model.Course
		notifications []model.Notification
		count         int
	)
	if err := ns.db.Find(&course, courseID).Error; err != nil {
		return nil, 0, err
	}
	ns.db.Model(&course).Offset(offset).Limit(limit).Association("Notifications").Find(&notifications)
	count = ns.db.Model(&course).Association("Notifications").Count()
	return notifications, count, nil
}
