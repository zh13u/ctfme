package routes

import (
	"ctfme/controllers"
	"ctfme/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", controllers.RegisterUser)
	api.Post("/login", controllers.LoginUser)
	api.Get("/challenges", middleware.JWTProtected(), controllers.GetChallenges)

	// setup
	api.Get("/admin/setup", middleware.JWTProtected(), controllers.GetSetup)
	api.Put("/admin/setup", middleware.JWTProtected(), controllers.UpdateSetup)

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

	// Route for admin to manage teams
	api.Get("/admin/teams", middleware.JWTProtected(), controllers.AdminGetTeams)
	api.Get("/admin/team/:id", middleware.JWTProtected(), controllers.AdminGetTeamDetail)

	//team features
	api.Post("/team/join", middleware.JWTProtected(), controllers.JoinTeam)
	api.Post("/team/leave", middleware.JWTProtected(), controllers.LeaveTeam)
	api.Post("/team/create", middleware.JWTProtected(), controllers.CreateTeam)

	api.Post("/submit", middleware.JWTProtected(), controllers.SubmitFlag)

	api.Get("/profile", middleware.JWTProtected(), controllers.GetProfile)

	api.Get("/scoreboard", controllers.GetScoreboard)
}
