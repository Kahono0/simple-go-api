package utils

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(
	// 	&models.User{},
	// 	&models.Item{},
	// 	&models.Order{},
	// )

	DB = db
}
