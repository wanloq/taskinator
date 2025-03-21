package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wanloq/taskinator/internal/controllers"
)

// SetupUserRoutes defines user-related routes
func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	// Protected route (requires authentication)
	userGroup.Get("/profile", controllers.GetUserProfile)
	userGroup.Put("/update", controllers.UpdateUserProfile)
	userGroup.Delete("/delete", controllers.DeleteUserProfile)
}
