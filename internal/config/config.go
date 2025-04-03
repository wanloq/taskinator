package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecretKey []byte

func ReadSecretFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadConfig loads environment variables
func LoadConfig() error {

	env := os.Getenv("ENV")
	if env == "" {
		log.Println("ENV is not set. Defaulting to 'development'.")
		env = "development"
	}

	// Load `.env` file in dev mode
	log.Println(env, " Environment")
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			return err
		}

		// Read JWT secret from environment variable
		JWTSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
		log.Println("Config successfully loaded for development environment")
		return nil
	}
	// load secrets from secrets dir

	// _, err := ReadSecretFile("./secrets/db_url")
	// if err != nil {
	// 	return fmt.Errorf("failed to load DB name: %w", err)
	// }

	// jwt_key, err := ReadSecretFile("./secrets/jwt_key")
	// if err != nil {
	// 	return fmt.Errorf("failed to load DB name: %w", err)
	// }

	// JWTSecretKey = []byte(jwt_key)

	// return fmt.Errorf("nothing loaded in %s environment", env)
	return nil
}
