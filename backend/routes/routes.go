package routes

import (
	"ctfme/controllers"
	"github.com/gofiber/fiber/v2"
	"ctfme/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", controllers.RegisterUser)
	api.Post("/login", controllers.LoginUser)
	api.Get("/challenges", controllers.GetChallenges)

	// Route for admin to manage challenges
	api.Post("/admin/challenge", middleware.JWTProtected(), controllers.CreateChallenge)
	api.Get("/admin/challenge", middleware.JWTProtected(), controllers.AdminGetChallenges)
	api.Put("/admin/challenge/:id", middleware.JWTProtected(), controllers.UpdateChallenge)
	api.Delete("/admin/challenge/:id", middleware.JWTProtected(), controllers.DeleteChallenge)

	// route for admin to manage users
	api.Get("/admin/users", middleware.JWTProtected(), controllers.AdminGetUser)
	api.Get("/admin/user/:id", middleware.JWTProtected(), controllers.AdminGetUserDetail)
	api.Put("/admin/user/:id", middleware.JWTProtected(), controllers.AdminUpdateUser)
	api.Delete("/admin/user/:id", middleware.JWTProtected(), controllers.AdminDeleteUser)
}
