package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

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

// @title Taskinator API
// @version 1.0
// @description A simple Task Manager API using Fiber and Swagger implemented in Go
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load config file
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
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

	// Routes
	routes.SetupRoutes(app)
	routes.SetupUserRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://0.0.0.0:%s/swagger/", port)
	if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
