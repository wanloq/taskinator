package utils

import (
	"errors"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
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
		userIDFloat, ok := (*claims)["user_id"].(float64) // JWT stores numbers as float64
		if !ok {
			return 0, errors.New("invalid token claims")
		}
		return uint(userIDFloat), nil
	}

	return 0, errors.New("invalid token")
}
