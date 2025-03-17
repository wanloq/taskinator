package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/wanloq/taskinator/docs"
	"github.com/wanloq/taskinator/internal/config"
	"github.com/wanloq/taskinator/internal/routes"
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

var tasks = []Task{
	{ID: 1, Name: "Learn Go", Status: false},
	{ID: 2, Name: "Learn Fiber", Status: false},
	{ID: 3, Name: "Integrate Swagger", Status: true},
}

type CustomClaims struct {
	UserID    uint   `json:"user_id"`
	UserEmail string `json:"user_email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

var secretKey = []byte("your_secret_key")

// @title Taskinator API
// @version 1.0
// @description A simple Task Manager API using Fiber and Swagger implemeneted in Go
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer config.CloseDB(db)

	// Run Migrations
	if err := config.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Server code
	app := fiber.New()
	app.Use(logger.New())

	// Register routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on http://localhost:%s/swagger/", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Home
func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to your Taskinator!")
}

// View all tasks
func GetAll(c *fiber.Ctx) error {
	if len(tasks) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No tasks found"})
	}
	return c.JSON(tasks)
}

// Create a new task
func CreateTask(c *fiber.Ctx) error {
	task := new(Task)
	err := c.BodyParser(task)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	id := len(tasks) + 1
	for _, task := range tasks {
		if task.ID >= id {
			id = task.ID + 1
		}
	}
	task.ID = id
	tasks = append(tasks, *task)
	return c.Status(fiber.StatusCreated).JSON(task)
}

// Mark a task as done
func FinishTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if len(tasks) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No tasks found"})
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = true
			return c.JSON(tasks[i])
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
}

// Delete a task
func DeleteTask(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if len(tasks) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No tasks found"})
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = true
			return c.JSON(tasks[i])
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
}

// // @Summary Get user profile
// // @Description Returns user profile if JWT is valid
// // @Tags Profile
// // @Security BearerAuth
// // @Produce json
// // @Success 200 {object} map[string]string "User profile"
// // @Failure 401 {object} map[string]string "Unauthorized"
// // @Router /profile [get]
// func Profile(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
// 	}

// 	tokenString := authHeader[len("Bearer "):]
// 	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	if err != nil || !token.Valid {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
// 	}

// 	claims, _ := token.Claims.(*CustomClaims)
// 	return c.JSON(fiber.Map{
// 		"user_id":   claims.UserID,
// 		"role":      claims.Role,
// 		"issued_at": claims.IssuedAt.Time.Format(time.RFC3339),
// 	})
// }

func GenerateJWT(userID uint, userEmail, role string) (string, error) {
	// var req TokenRequest
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)

}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		fmt.Println("User ID:", claims.UserID)
		fmt.Println("Role:", claims.Role)
		return token, nil
	}

	return nil, fmt.Errorf("invalid token")
}
