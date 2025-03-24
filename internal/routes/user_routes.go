package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wanloq/taskinator/internal/controllers"
	"github.com/wanloq/taskinator/internal/middleware"
)

// SetupUserRoutes defines user-related routes
func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	// Protected route (requires authentication)
	userGroup.Get("/profile", controllers.GetUserProfile)
	userGroup.Put("/update", controllers.UpdateUserProfile)
	userGroup.Post("/password-reset/request", controllers.RequestPasswordReset)
	userGroup.Post("/password-reset/confirm", controllers.PasswordReset)
	userGroup.Post("/email/verify/request", controllers.RequestEmailVerification)
	userGroup.Get("/email/verify", controllers.VerifyEmail)
	// userGroup.Delete("/delete", controllers.DeleteUserProfile)
	userGroup.Delete("/admin/delete-user/:id", middleware.RoleMiddleware("admin"), controllers.DeleteUserProfile)
}
