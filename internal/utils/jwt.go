package utils

import (
	"errors"
	"time"

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
