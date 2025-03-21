package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wanloq/taskinator/internal/dto"
	"github.com/wanloq/taskinator/internal/repositories"
	"github.com/wanloq/taskinator/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

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
