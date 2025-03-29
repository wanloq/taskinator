package controllers

import (
	"log"
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

	// Check if user already exists
	existingUser, _ := repositories.GetUserByEmail(req.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already in use"})
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// user model
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsVerified:   false,
	}
	user.PasswordHash = hashedPassword

	// Save user to database
	if err := repositories.CreateUser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate verification token
	verificationToken, err := utils.GenerateEmailVerificationToken(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate verification token"})
	}

	// Send verification email
	log.Println("Sending verification mail to ", user.Email)
	go utils.SendVerificationEmail(user.Email, verificationToken)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "User registered successfully. Please verify your email."})
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

	if !user.IsVerified {
		go func() {
			// Generate email verification token
			verificationToken, err := utils.GenerateEmailVerificationToken(user.Email)
			if err != nil {
				log.Println("Could not generate verification token", err)
				return
			}

			// Send verification email
			log.Println("Sending verification mail to ", user.Email)
			err = utils.SendVerificationEmail(user.Email, verificationToken)
			if err != nil {
				log.Println("Could not send verification mail to ", user.Email, err)
				return
			}
		}()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User email not verified"})
	}
	// Compare hashed password
	if !utils.ComparePasswords(user.PasswordHash, req.Password) {
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
	// Extract user ID from JWT token
	userID, _, err := utils.ExtractUserFromToken(c)
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
		"status":   user.IsVerified,
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
	userID, _, err := utils.ExtractUserFromToken(c)
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
		user.PasswordHash = string(hashedPassword)
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
// // @Router /user/delete [delete]
// @Router /user/admin/delete-user/:id [delete]
func DeleteUserProfile(c *fiber.Ctx) error {
	// Extract user ID from JWT token
	userID, _, err := utils.ExtractUserFromToken(c)
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

// @Summary Request Email Verification
// @Description Sends an email with a verification link to the user
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param request body dto.RequestEmailVerification true "Email verification request"
// @Success 200 {object} map[string]string "Verification email sent successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/email/verify/request [post]
func RequestEmailVerification(c *fiber.Ctx) error {
	var req dto.RequestEmailVerification
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Check if user exists
	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate email verification token
	verificationToken, err := utils.GenerateEmailVerificationToken(user.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate verification token"})
	}

	// Send verification email
	log.Println("Sending verification mail to ", user.Email)
	err = utils.SendVerificationEmail(user.Email, verificationToken)
	if err != nil {
		log.Println("Could not Send verification mail to ", user.Email, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Verification email sent successfully"})
}

// @Summary Verify Email
// @Description Confirms the user's email verification
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param token query string true "Verification Token"
// @Success 200 {object} map[string]string "Email verified successfully"
// @Failure 400 {object} map[string]string "Invalid or expired token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/email/verify [get]
func VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing verification token"})
	}

	email, err := utils.VerifyEmailToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Extract user info
	user, err := repositories.GetUserByEmail(email)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if user.IsVerified {
		return c.JSON(fiber.Map{"message": "Email already verified."})
	}

	// Update and Save user
	user.IsVerified = true
	if err := repositories.UpdateUser(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update user"})
	}

	return c.JSON(fiber.Map{"message": "Email successfully verified. You may now log in."})
}

// @Summary Request Password Reset
// @Description Sends a password reset link to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.PasswordResetRequest true "User email for password reset"
// @Success 200 {object} map[string]string "Password reset link sent"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/password-reset/request [post]
func RequestPasswordReset(c *fiber.Ctx) error {
	var req dto.PasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := repositories.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Generate password reset token
	resetToken, err := utils.GeneratePasswordResetToken(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate reset token"})
	}

	// Send email with password reset link
	go utils.SendPasswordResetEmail(user.Email, resetToken)

	return c.JSON(fiber.Map{"message": "Password reset link sent to your email."})
}

// @Summary Reset User Password
// @Description Verifies reset token and updates user password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResetRequest true "New password and token"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid or expired token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/password-reset/confirm [post]
func PasswordReset(c *fiber.Ctx) error {
	var req dto.ResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	email, err := utils.VerifyPasswordResetToken(req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	err = repositories.UpdateUserPassword(email, string(hashedPassword))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update password"})
	}

	return c.JSON(fiber.Map{"message": "Password successfully reset."})
}
