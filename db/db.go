package db

import (
	"CourseService/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func New() *gorm.DB {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_NAME"),
		))
	if err != nil {
		fmt.Print("Storage error: ", err)
	}
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	return db
}

func TestDB() *gorm.DB {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_PASSWORD"),
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("TEST_DATABASE_NAME"),
		))
	if err != nil {
		fmt.Print("Storage error: ", err)
	}
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(10)
	db.LogMode(false)
	return db
}

func DropTestDB(db *gorm.DB) error {
	return db.Exec(fmt.Sprint("DROP DATABASE ", os.Getenv("TEST_DATABASE_NAME"))).Error
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Course{},
		&model.Review{},
		&model.Notification{},
		&model.Student{},
	)
}
