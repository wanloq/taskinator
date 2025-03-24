package utils

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wanloq/taskinator/internal/config"
	"github.com/wanloq/taskinator/internal/dto"
)

// GenerateJWT creates a JWT token
func GenerateJWT(userID uint, email, role string) (string, error) {
	claims := dto.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

// VerifyJWT verifies and extracts claims from a token
func VerifyJWT(tokenString string) (*dto.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return config.JWTSecretKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*dto.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// ExtractUserIDFromToken extracts the user_id from the JWT token
func ExtractUserIDFromToken(c *fiber.Ctx) (uint, error) {
	// Get the token from the request header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing token")
	}

	// Extract the token (remove "Bearer " prefix)
	tokenString := authHeader[7:]

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecretKey, nil
	})
	if err != nil {
		return 0, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return uint(userIDFloat), nil
	}

	return 0, errors.New("invalid token")
}

func ExtractUserFromToken(c *fiber.Ctx) (uint, string, error) {
	// Get the token from the request header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, "", errors.New("missing token")
	}

	// Extract the token (remove "Bearer " prefix)
	tokenString := authHeader[7:]

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return config.JWTSecretKey, nil
	})
	if err != nil {
		return 0, "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		userID := uint((*claims)["user_id"].(float64))
		role := (*claims)["role"].(string)
		return userID, role, nil
	}

	return 0, "", errors.New("invalid token")
}

func GenerateEmailVerificationToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

func VerifyEmailToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return config.JWTSecretKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return (*claims)["email"].(string), nil
	}

	return "", errors.New("invalid token")
}

// GeneratePasswordResetToken creates a JWT token for password reset
func GeneratePasswordResetToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

// VerifyPasswordResetToken verifies the JWT password reset token
func VerifyPasswordResetToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	email, exists := claims["email"].(string)
	if !exists {
		return "", errors.New("invalid token claims")
	}
	log.Println("Password reset token Verified")
	return email, nil
}
