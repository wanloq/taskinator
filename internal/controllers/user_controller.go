package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wanloq/taskinator/internal/dto"
	"github.com/wanloq/taskinator/internal/models"
	"github.com/wanloq/taskinator/internal/repositories"
	"github.com/wanloq/taskinator/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// @Summary User Registration
// @Description RegisterUser handles user registration: Creates a new user and returns a success message or an error.
// @Tags Registration
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration request"
// @Success 201 {object} map[string]string "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/register [post]
func RegisterUser(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword := req.Password

	// user model
	user := models.User{
		Username:      req.Username,
		Email:         req.Email,
		Password_Hash: hashedPassword,
	}
	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user.Password_Hash = hashedPassword

	// Save user to database
	if err := repositories.CreateUser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// @Summary User Login
// @Description LoginUser handles user authentication: Logs in a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string "Token response"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/login [post]
func LoginUser(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Fetch user from database
	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// Compare hashed password
	if !utils.ComparePasswords(user.Password_Hash, req.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token (we already implemented this)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})
}

// @Summary Get user profile
// @Description Returns the currently logged-in user's profile if JWT is valid
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /user/profile [get]
func GetUserProfile(c *fiber.Ctx) error {

	// userID := c.Locals("user_id")
	// email := c.Locals("email")
	// role := c.Locals("role")

	// Extract user ID from JWT token
	userID, err := utils.ExtractUserIDFromToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fetch user from the database
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Return user details
	return c.JSON(fiber.Map{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// @Summary Update user profile
// @Description UpdateUserProfile updates the authenticated user's profile if JWT is valid
// @Tags Update Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateRequest true "Update Request"
// @Success 200 {object} map[string]string "Profile Updated successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User Not Found"
// @Router /user/update [put]
func UpdateUserProfile(c *fiber.Ctx) error {
	// Extract user ID from JWT token
	userID, err := utils.ExtractUserIDFromToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var req dto.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Fetch the existing user
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the new email is already in use by another user
	existingUser, err := repositories.GetUserByEmail(req.Email)
	if err == nil && existingUser.ID != user.ID {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Email already in use"})
	}

	// Update user fields
	user.Username = req.Username
	user.Email = req.Email

	// If password is provided, hash and update it
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
		}
		user.Password_Hash = string(hashedPassword)
	}

	// Save updated user
	if err := repositories.UpdateUser(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}

	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}

// @Summary Delete user profile
// @Description DeleteUserProfile Removes the authenticated user's profile if JWT is valid
// @Tags Delete Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Profile deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User Not Found"
// @Router /user/delete [delete]
func DeleteUserProfile(c *fiber.Ctx) error {
	// Extract user ID from JWT token
	userID, err := utils.ExtractUserIDFromToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fetch the existing user
	user, err := repositories.GetUserByID(userID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Delete user
	if err := repositories.DeleteUser(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete user"})
	}

	return c.JSON(fiber.Map{"message": "Profile deleted successfully"})
}
