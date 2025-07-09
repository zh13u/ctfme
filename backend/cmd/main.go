package main

import (
	"ctfme/config"
	"ctfme/database"
	"ctfme/database/migrations"
	"ctfme/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	database.ConnectDB()

	// Run migration to remove config fields from challenges table
	migrations.RemoveChallengeConfigFields()

	// Thêm migration thêm trường difficulty
	migrations.AddDifficultyToChallenge()

	migrations.AddCurrentPointsToChallenge()

	config.InitConfig()  // Initialize configuration after .env is loaded
	config.PrintConfig() // Print configuration for debugging
	app := fiber.New()

	// Thêm middleware CORS ngay sau khi khởi tạo app
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))

}
