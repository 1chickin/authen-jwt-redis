package config

import (
	"fmt"
	"os"

	"github.com/1chickin/authen-jwt-redis/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := "disable"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// DB, err = gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func MigrateDB() {
	DB.AutoMigrate(&model.User{})
}
