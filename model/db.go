package model

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define the Letters struct which represents the Letters table
type Letters struct {
	ID       uint   
	From     string `json:"from"`
	FromName string `json:"from_name"`
	ToUser   string `json:"to"`
	Message  string `json:"message"`
}

// Retrieve environment variables
var dbHost string
var db *gorm.DB

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}
	dbHost = os.Getenv("DB_URL")

	db, err = gorm.Open(postgres.Open(dbHost), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}
	// Auto migrate the User table
	// err = db.AutoMigrate(&Letters{})
	if err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}
}

func AddLetter(letter *Letters) *Letters {
	// Create a new User
	user := letter
	db.Create(&user)
	if db.Error != nil {
		log.Fatal("Failed to create user:", db.Error)
	}
	return user
}

func GetUserMessage(to string) []Letters {
	var letters []Letters
	result := db.Where(`to_user = ?`, to).Find(&letters)
	if result.Error != nil {
		log.Fatal("Failed to retrieve user:", result.Error)
	}
	return letters
}
