package model

import "github.com/jinzhu/gorm"

type Course struct {
	gorm.Model
	Mentor        uint   `gorm:"not null"`
	Name          string `gorm:"type:varchar(100);not null"`
	Description   string `gorm:"type:varchar(300);not null"`
	Prerequisite  string `gorm:"type:varchar(300);not null"`
	Aim           string `gorm:"type:varchar(300);not null"`
	MaxStudent    uint   `gorm:"not null"`
	Fee           uint   `gorm:"not null"`
	Rating        uint   `gorm:"not null"`
	Reviews       []Review
	Notifications []Notification
	Students      []Student `gorm:"many2many:student_course"`
}

type Student struct {
	ID      uint
	Courses []Course `gorm:"many2many:student_course"`
}