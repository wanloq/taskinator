package routes

import (
	"github.com/wanloq/taskinator/internal/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Public routes
	app.Get("/swagger/*", swagger.HandlerDefault)
	api.Post("/register", controllers.RegisterUser)
	api.Post("/login", controllers.LoginUser)

	// Task routes
	// app.Get("/all", GetAll)
	// app.Post("/task", CreateTask)
	// app.Patch("/task", FinishTask)
}
