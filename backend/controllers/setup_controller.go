package controllers

import (
    "ctfme/database"
    "ctfme/models"
    "github.com/gofiber/fiber/v2"
)

func GetSetup(c *fiber.Ctx) error {
    var config models.Setup
    if err := database.DB.First(&config).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Config not found"})
    }
    return c.JSON(config)
}

func UpdateSetup(c *fiber.Ctx) error {
    var input models.Setup
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }
    var config models.Setup
    if err := database.DB.First(&config).Error; err != nil {
        // Nếu chưa có thì tạo mới
        config = input
        if err := database.DB.Create(&config).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Could not create config"})
        }
    } else {
        // Cập nhật
        config.CTFMode = input.CTFMode
        config.DynamicScoreEnabled = input.DynamicScoreEnabled
        config.DynamicScoreDecay = input.DynamicScoreDecay
        config.DynamicScoreMin = input.DynamicScoreMin
        if err := database.DB.Save(&config).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Could not update config"})
        }
    }
    return c.JSON(config)
}