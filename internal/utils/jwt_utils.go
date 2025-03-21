package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/wanloq/taskinator/internal/config"
)

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
		userIDFloat, ok := (*claims)["user_id"].(float64) // JWT stores numbers as float64
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return uint(userIDFloat), nil
	}

	return 0, errors.New("invalid token")
}
