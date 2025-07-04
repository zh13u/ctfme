package controllers

import (
	"ctfme/models"
	"github.com/gofiber/fiber/v2"
	"ctfme/database"
)

func GetChallenges(c *fiber.Ctx) error {
	var challenges []models.Challenge

	// get all challenge visible to user
	if err := database.DB.Where("visible = ?", true).Find(&challenges).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch challenges"})
	}

	// hide flag
	type ChallengeResponse struct {
		ID			uint	`json:"id"`
		Title		string	`json:"title"`
		Description	string	`json:"description"`
		Category	string	`json:"category"`
		Points		int		`json:"points"`	
		FileURL		string	`json:"file_url"`
	}

	var response []ChallengeResponse
	for _, challenge := range challenges {
		response = append(response, ChallengeResponse{
			ID: challenge.ID,
			Title: challenge.Title,
			Description: challenge.Description,
			Category: challenge.Category,
			Points: challenge.Points,
			FileURL: challenge.FileURL,
		})
	}

	return c.JSON(response)
}

// API for admin to create new challenge
func CreateChallenge(c *fiber.Ctx) error {
	// Lấy user_id từ JWT middleware
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}

	type ChallengeInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Points      int    `json:"points"`
		Flag        string `json:"flag"`
		FileURL     string `json:"file_url"`
		Visible     bool   `json:"visible"`
	}
	var input ChallengeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	challenge := models.Challenge{
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		Points:      input.Points,
		Flag:        input.Flag,
		FileURL:     input.FileURL,
		Visible:     input.Visible,
	}
	if err := database.DB.Create(&challenge).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create challenge"})
	}
	return c.JSON(fiber.Map{"message": "Challenge created!", "challenge": challenge})
}

// API cho admin lấy danh sách tất cả challenge (bao gồm cả ẩn và flag)
func AdminGetChallenges(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	var challenges []models.Challenge
	if err := database.DB.Find(&challenges).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch challenges"})
	}
	return c.JSON(challenges)
}

// API cho admin sửa challenge
func UpdateChallenge(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	id := c.Params("id")
	var challenge models.Challenge
	if err := database.DB.First(&challenge, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Challenge not found"})
	}
	type ChallengeInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Points      int    `json:"points"`
		Flag        string `json:"flag"`
		FileURL     string `json:"file_url"`
		Visible     bool   `json:"visible"`
	}
	var input ChallengeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	challenge.Title = input.Title
	challenge.Description = input.Description
	challenge.Category = input.Category
	challenge.Points = input.Points
	challenge.Flag = input.Flag
	challenge.FileURL = input.FileURL
	challenge.Visible = input.Visible
	if err := database.DB.Save(&challenge).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update challenge"})
	}
	return c.JSON(fiber.Map{"message": "Challenge updated!", "challenge": challenge})
}

// API cho admin xóa challenge
func DeleteChallenge(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	if !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	id := c.Params("id")
	if err := database.DB.Unscoped().Delete(&models.Challenge{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete challenge"})
	}
	return c.JSON(fiber.Map{"message": "Challenge deleted!"})
}