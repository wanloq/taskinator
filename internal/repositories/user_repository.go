package repositories

import (
	"github.com/wanloq/taskinator/internal/config"
	"github.com/wanloq/taskinator/internal/models"
)

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) error {
	result := config.DB.Create(user)
	return result.Error
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

// UpdateUser updates an existing user in the database
func DeleteUser(user *models.User) error {
	return config.DB.Delete(user).Error
}

// UpdateUserPassword updates a user's password by email
func UpdateUserPassword(email string, newPasswordHash string) error {
	result := config.DB.Model(&models.User{}).Where("email = ?", email).Update("password_hash", newPasswordHash)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// VerifyUserEmail sets the 'is_verified' field to true
func VerifyUserEmail(email string) error {
	result := config.DB.Model(&models.User{}).Where("email = ?", email).Update("is_verified", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
