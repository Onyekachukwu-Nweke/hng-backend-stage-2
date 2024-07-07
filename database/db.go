package database

import (
	// "fmt"
	"log"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(connectionString string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB.AutoMigrate(&models.User{}, &models.Organisation{})
}

func GetDB() *gorm.DB {
	return DB
}
