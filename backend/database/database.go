package database

import (
	"log"
	"os"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"ctfme/models"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func ConnectDB() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB =db

	db.AutoMigrate(&models.Team{}, &models.User{}, &models.Challenge{}, &models.Submission{})
}