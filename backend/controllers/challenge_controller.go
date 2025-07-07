package controllers

import (
	"ctfme/config"
	"ctfme/database"
	"ctfme/models"

	"github.com/gofiber/fiber/v2"
)

// calculateDynamicPoints calculates the current points for a challenge based on solve count
func calculateDynamicPoints(challenge models.Challenge, solveCount int64) int {
	points := challenge.Points
	if config.DynamicScoreEnabled && solveCount > 0 {
		if config.DynamicScoreDecay > 0 {
			points = challenge.Points - int(solveCount-1)*config.DynamicScoreDecay
			if points < config.DynamicScoreMin {
				points = config.DynamicScoreMin
			}
		}
	}
	return points
}

func GetChallenges(c *fiber.Ctx) error {
	var challenges []models.Challenge

	// get all challenge visible to user
	if err := database.DB.Where("visible = ?", true).Find(&challenges).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch challenges"})
	}

	// hide flag
	type ChallengeResponse struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Points      int    `json:"points"`
		FileURL     string `json:"file_url"`
		SolveCount  int64  `json:"solve_count"`
	}

	var response []ChallengeResponse
	for _, challenge := range challenges {
		var solveCount int64
		if config.TeamMode {
			database.DB.Raw(`
				SELECT COUNT(DISTINCT u.team_id)
				FROM submissions s
				JOIN users u ON s.user_id = u.id
				WHERE s.challenge_id = ? AND s.is_correct = true AND u.team_id IS NOT NULL
			`, challenge.ID).Scan(&solveCount)
		} else {
			database.DB.Model(&models.Submission{}).
				Where("challenge_id = ? AND is_correct = true", challenge.ID).
				Distinct("user_id").Count(&solveCount)
		}
		// Calculate dynamic points
		displayPoints := calculateDynamicPoints(challenge, solveCount)

		response = append(response, ChallengeResponse{
			ID:          challenge.ID,
			Title:       challenge.Title,
			Description: challenge.Description,
			Category:    challenge.Category,
			Points:      displayPoints,
			FileURL:     challenge.FileURL,
			SolveCount:  solveCount,
		})
	}

	return c.JSON(response)
}

// API for admin to create new challenge
func CreateChallenge(c *fiber.Ctx) error {
	// get user_id from JWT middleware
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

// API for admin to get all challenges
func AdminGetChallenges(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	var challenges []models.Challenge
	database.DB.Preload("SolvedBy").Find(&challenges)
	return c.JSON(challenges)
}

// API for admin to update challenge
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

// API for admin to delete challenge
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

	// First, find the challenge
	var challenge models.Challenge
	if err := database.DB.First(&challenge, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Challenge not found"})
	}

	// Delete related submissions first
	if err := database.DB.Where("challenge_id = ?", id).Delete(&models.Submission{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete challenge submissions"})
	}

	// Clear the many-to-many relationship
	if err := database.DB.Model(&challenge).Association("SolvedBy").Clear(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not clear challenge solvers"})
	}

	// Now delete the challenge
	if err := database.DB.Unscoped().Delete(&challenge).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete challenge"})
	}

	return c.JSON(fiber.Map{"message": "Challenge deleted!"})
}

func SubmitFlag(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Get user info
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if user is in a team when in team mode
	teamWarning := ""
	if config.TeamMode && user.TeamID == nil {
		teamWarning = "Warning: You are not in a team. Consider joining a team for better experience in team mode."
	}

	type SubmitInput struct {
		ChallengeID uint   `json:"challenge_id"`
		Flag        string `json:"flag"`
	}
	var input SubmitInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	var challenge models.Challenge
	if err := database.DB.First(&challenge, input.ChallengeID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Challenge not found"})
	}
	isCorrect := input.Flag == challenge.Flag

	// Check if user/team already solved this challenge
	var existingSolve bool
	var debugSolveCount int64

	if config.TeamMode {
		if user.TeamID != nil {
			// In team mode with team, check if team already solved
			database.DB.Raw(`
				SELECT COUNT(*)
				FROM submissions s
				JOIN users u ON s.user_id = u.id
				WHERE s.challenge_id = ? AND s.is_correct = true AND u.team_id = ?
			`, challenge.ID, user.TeamID).Scan(&debugSolveCount)
			existingSolve = debugSolveCount > 0
		} else {
			// In team mode without team, check if user already solved
			database.DB.Model(&models.Submission{}).
				Where("challenge_id = ? AND user_id = ? AND is_correct = true", challenge.ID, userID).
				Count(&debugSolveCount)
			existingSolve = debugSolveCount > 0
		}
	} else {
		// In user mode, check if user already solved
		database.DB.Model(&models.Submission{}).
			Where("challenge_id = ? AND user_id = ? AND is_correct = true", challenge.ID, userID).
			Count(&debugSolveCount)
		existingSolve = debugSolveCount > 0
	}

	// Debug: Print solve status
	println("=== SubmitFlag Debug ===")
	println("UserID:", userID)
	println("TeamMode:", config.TeamMode)
	println("User.TeamID:", user.TeamID)
	println("ExistingSolve:", existingSolve)
	println("DebugSolveCount:", debugSolveCount)
	println("========================")

	// save submission
	sub := models.Submission{
		UserID:      userID,
		TeamID:      user.TeamID, // Include team ID in submission
		ChallengeID: challenge.ID,
		Flag:        input.Flag,
		IsCorrect:   isCorrect,
	}
	database.DB.Create(&sub)

	if isCorrect {
		if !existingSolve {
			// Only add to SolvedBy if this is the first solve for user/team
			database.DB.Model(&challenge).Association("SolvedBy").Append(&user)
		}
		// calculate points
		var solveCount int64
		if config.TeamMode {
			database.DB.Raw(`
				SELECT COUNT(DISTINCT u.team_id)
				FROM submissions s
				JOIN users u ON s.user_id = u.id
				WHERE s.challenge_id = ? AND s.is_correct = true AND u.team_id IS NOT NULL
			`, challenge.ID).Scan(&solveCount)
		} else {
			database.DB.Model(&models.Submission{}).
				Where("challenge_id = ? AND is_correct = true", challenge.ID).
				Count(&solveCount)
		}
		points := calculateDynamicPoints(challenge, solveCount)
		response := fiber.Map{"result": "Correct!", "points": points}
		if teamWarning != "" {
			response["warning"] = teamWarning
		}
		if existingSolve {
			response["message"] = "You already solved this challenge before"
		}
		return c.JSON(response)
	}

	response := fiber.Map{"result": "Incorrect!"}
	if teamWarning != "" {
		response["warning"] = teamWarning
	}
	return c.JSON(response)
}
