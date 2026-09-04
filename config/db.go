package config

import (
	"log"
	"os"

	"example.com/event-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("Environtment variable not found")
		return
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection failured")
	}

	err = database.AutoMigrate(&models.Event{}, &models.User{}, &models.Booking{})
	if err != nil {
		log.Fatal("Database migration failed :", err)
	}

	DB = database
	log.Println("Database Connected Successfully")
}