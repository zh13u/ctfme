package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "ctfme/database"
	"ctfme/routes"
	"github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    database.ConnectDB()
    app := fiber.New()
	routes.SetupRoutes(app)
    log.Fatal(app.Listen(":3000"))
}