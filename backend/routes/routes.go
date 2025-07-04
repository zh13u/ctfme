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

	// Route cho admin tạo challenge
	api.Post("/admin/challenge", middleware.JWTProtected(), controllers.CreateChallenge)

	// Route cho admin quản lý challenge
	api.Get("/admin/challenge", middleware.JWTProtected(), controllers.AdminGetChallenges)
	api.Put("/admin/challenge/:id", middleware.JWTProtected(), controllers.UpdateChallenge)
	api.Delete("/admin/challenge/:id", middleware.JWTProtected(), controllers.DeleteChallenge)

	// need login
	
}