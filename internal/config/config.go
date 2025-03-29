package config

import (
	"os"

	"github.com/joho/godotenv"
)

var JWTSecretKey []byte

// LoadConfig loads environment variables
func LoadConfig() error {
	// Load `.env` file
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Read JWT secret from environment variable
	JWTSecretKey = []byte(os.Getenv("JWT_SECRET"))
	return nil
}
