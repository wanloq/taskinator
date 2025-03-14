package repositories

import (
	"github.com/wanloq/taskinator/internal/db"
	"github.com/wanloq/taskinator/internal/models"
)

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) error {
	result := db.DB.Create(user)
	return result.Error
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
