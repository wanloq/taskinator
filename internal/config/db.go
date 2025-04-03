package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes and connects to the PostgreSQL database
func ConnectDB() (*gorm.DB, error) {

	// Load DatabaseURL from secrets file
	// databaseURL, err := LoadDBConfig()
	// if err != nil {

	// 	log.Println(err)
	// 	return nil, err
	// }

	// Prioritize DATABASE_URL if set (for Docker compatibility)
	databaseURL, err := LoadDBURL()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	// if databaseURL == "" {
	// 	databaseURL = fmt.Sprintf(
	// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 		os.Getenv("DB_HOST"),
	// 		os.Getenv("DB_USER"),
	// 		os.Getenv("DB_PASSWORD"),
	// 		os.Getenv("DB_NAME"),
	// 		os.Getenv("DB_PORT"),
	// 	)
	// }

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB = db
	fmt.Println("✅ Connected to the database")
	return DB, nil
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error retrieving database connection:", err)
		return
	}
	sqlDB.Close()
}

// LoadEmailConfig returns DB configurations (Username, Password, Name and Host) or an error .
func LoadDBURL() (DBURL string, err error) {

	DBURL, err = ReadSecretFile("./secrets/db_url")
	if err != nil {
		return "", fmt.Errorf("failed to load DB_URL: %w", err)
	}

	log.Println("Database URL:", DBURL)
	return
}
