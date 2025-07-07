package main

import (
	"ctfme/config"
	"ctfme/database"
	"ctfme/routes"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	config.InitConfig()  // Initialize configuration after .env is loaded
	config.PrintConfig() // Print configuration for debugging
	database.ConnectDB()
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
