package db

import (
	"log"
	"restful-service/models"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDb *gorm.DB

func Init(dataSourceName string) error {
	var err error

	dsn := strings.Replace(dataSourceName, "TaskManager", "todo", 1)
	gormDb, err = gorm.Open(mysql.Open(dsn + "?parseTime=true"))
	log.Println("Database connection established")

	if err != nil {
		log.Println(err)
	}

	gormDb.AutoMigrate(&models.User{})
	gormDb.AutoMigrate(&models.Task{})
	log.Println("Database migrated")

	return nil
}

func GetGormDb() *gorm.DB {
	return gormDb
}

func Close() {
	db, _ := gormDb.DB()
	db.Close()
	log.Println("Database connection closed")
}
