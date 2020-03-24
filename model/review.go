package model

import "github.com/jinzhu/gorm"

type Review struct {
	gorm.Model
	User     uint   `gorm:"not null"`
	Star     uint   `gorm:"not null"`
	Content  string `gorm:"type:varchar(300);not null"`
	CourseID uint
	Course   Course
}
