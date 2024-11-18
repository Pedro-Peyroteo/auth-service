package database

import (
	"auth-service/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Instance *gorm.DB
var dbError error

func Connect() {
	// Load environment variables from .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Construct the connection string from environment variables
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Create a custom logger with Silent log level
	customLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Set log level to Silent
			Colorful:      true,          // Optional: Colorful output
		},
	)

	// Connect to the database
	Instance, dbError = gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: customLogger})

	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to Database!")
	}

	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
	log.Println("Database Migration Completed!")
}
