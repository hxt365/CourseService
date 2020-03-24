package model

import (
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	Content  string `gorm:"type:varchar(200);not null"`
	CourseID uint
	Course   Course
}
