package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wanloq/taskinator/internal/handlers"
)

// SetupUserRoutes defines user-related routes
func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	// Protected route (requires authentication)
	userGroup.Get("/profile", handlers.GetUserProfile)
	userGroup.Put("/update", handlers.UpdateUserProfile)
	userGroup.Delete("/delete", handlers.DeleteUserProfile)
}
