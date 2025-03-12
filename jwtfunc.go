package main

// // var secretKey = []byte("your_secret_key")
// type CustomClaims struct {
// 	UserID string `json:"user_id"`
// 	Role   string `json:"role"`
// 	jwt.RegisteredClaims
// }

// func GenerateJWT() (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": "wan",                                             // What should go here?
// 		"role":    "admin",                                           // What should go here?
// 		"exp":     jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // How do we set this to 1 hour from now?
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(secretKey)

// }
