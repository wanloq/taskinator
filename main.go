package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/wanloq/taskinator/docs"
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
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Dummy user credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var secretKey = []byte("your_secret_key")

// @title Taskinator API
// @version 1.0
// @description A simple Task Manager API using Fiber and Swagger implemeneted in Go
// @host localhost:3000
// @BasePath /
// @schemes http
func main() {

	// // Generate JWT
	// tokenString, err := GenerateJWT()
	// if err != nil {
	// 	fmt.Println("Error generating token:", err)
	// 	return
	// }
	// fmt.Println("\nGenerated Token:", tokenString)
	// // tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDE2MzM0NTEsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOiJ3YW4ifQ.cg7wNVP_-RZgbJy5_EQusUtinWE2wtd09t71oUmc5mc"

	// // Verify JWT
	// token, err := VerifyJWT(tokenString)
	// if err != nil {
	// 	fmt.Println("Error verifying token:", err)
	// 	return
	// }
	// fmt.Println("\nVerified Token:", token)

	// Server code
	app := fiber.New()
	app.Use(logger.New())

	// Routes
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Login Route
	app.Post("/login", Login)
	app.Get("/all", GetAll)
	app.Post("/task", CreateTask)
	app.Patch("/task", FinishTask)

	log.Println("\nServer running @ http://localhost:3000/swagger/")
	log.Fatal(app.Listen(":3000"))
	// Handlers
}

// Home
func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to your Taskinator!")
}

// @Summary User Login
// @Description Logs in a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} map[string]string "Token response"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /login [post]
func Login(c *fiber.Ctx) error {

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Dummy validation (replace with DB check)
	if req.Username != "admin" || req.Password != "password" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT Token
	token, err := GenerateJWT(req.Username, "admin")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
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

func GenerateJWT(userID, role string) (string, error) {
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
