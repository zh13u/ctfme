package main

import (
	"ctfme/config"
	"ctfme/database"
	"ctfme/database/migrations"
	"ctfme/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	database.ConnectDB()

	// Run migration to remove config fields from challenges table
	migrations.RemoveChallengeConfigFields()

	config.InitConfig()  // Initialize configuration after .env is loaded
	config.PrintConfig() // Print configuration for debugging
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
