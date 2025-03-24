package models

import "gorm.io/gorm"

// User represents the users table
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:user"`
	IsVerified   bool   `gorm:"default:false"`
}
