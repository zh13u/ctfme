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
	if err := database.DB.Preload("SolvedBy").Where("visible = ?", true).Find(&challenges).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Không thể lấy danh sách thử thách"})
	}
	// Ẩn trường flag trước khi trả về
	type ChallengePublic struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Points      int    `json:"points"`
		FileURL     string `json:"fileURL"`
		Visible     bool   `json:"visible"`
		Difficulty  string `json:"difficulty"`
		// Không có trường flag
	}
	var publicChallenges []ChallengePublic
	for _, ch := range challenges {
		// solveCount := int64(len(ch.SolvedBy))
		// currentPoints := calculateDynamicPoints(ch, solveCount)
		publicChallenges = append(publicChallenges, ChallengePublic{
			ID:          ch.ID,
			Title:       ch.Title,
			Description: ch.Description,
			Category:    ch.Category,
			Points:      ch.CurrentPoints,
			FileURL:     ch.FileURL,
			Visible:     ch.Visible,
			Difficulty:  ch.Difficulty,
		})
	}
	return c.JSON(publicChallenges)
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
		Difficulty  string `json:"difficulty"`
	}
	var input ChallengeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	challenge := models.Challenge{
		Title:         input.Title,
		Description:   input.Description,
		Category:      input.Category,
		Points:        input.Points,
		Flag:          input.Flag,
		FileURL:       input.FileURL,
		Visible:       input.Visible,
		Difficulty:    input.Difficulty,
		CurrentPoints: input.Points,
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

	// Clean response
	type UserResponse struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		TeamID   *uint  `json:"team_id"`
	}
	type ChallengeResponse struct {
		ID           uint           `json:"id"`
		CreatedAt    string         `json:"created_at"`
		UpdatedAt    string         `json:"updated_at"`
		DeletedAt    *string        `json:"deleted_at"`
		Title        string         `json:"title"`
		Description  string         `json:"description"`
		Category     string         `json:"category"`
		Points       int            `json:"points"`
		Flag         string         `json:"flag"`
		FileURL      string         `json:"file_url"`
		Visible      bool           `json:"visible"`
		SolvedBy     []UserResponse `json:"solved_by"`
		DynamicScore bool           `json:"dynamic_score"`
		MinScore     int            `json:"min_score"`
		Decay        int            `json:"decay"`
	}

	// Get system configuration
	var setup models.Setup
	if err := database.DB.First(&setup).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch system configuration"})
	}

	var response []ChallengeResponse
	for _, ch := range challenges {
		var solvedBy []UserResponse
		for _, u := range ch.SolvedBy {
			solvedBy = append(solvedBy, UserResponse{
				ID:       u.ID,
				Username: u.Username,
				Email:    u.Email,
				TeamID:   u.TeamID,
			})
		}
		var deletedAt *string
		if ch.DeletedAt.Valid {
			dt := ch.DeletedAt.Time.Format("2006-01-02T15:04:05Z07:00")
			deletedAt = &dt
		}
		response = append(response, ChallengeResponse{
			ID:           ch.ID,
			CreatedAt:    ch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    ch.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			DeletedAt:    deletedAt,
			Title:        ch.Title,
			Description:  ch.Description,
			Category:     ch.Category,
			Points:       ch.Points,
			Flag:         ch.Flag,
			FileURL:      ch.FileURL,
			Visible:      ch.Visible,
			SolvedBy:     solvedBy,
			DynamicScore: setup.DynamicScoreEnabled,
			MinScore:     setup.DynamicScoreMin,
			Decay:        setup.DynamicScoreDecay,
		})
	}
	return c.JSON(response)
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

	// Chặn submit nếu chưa vào team ở team mode
	if config.TeamMode && user.TeamID == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bạn cần vào team mới được submit!"})
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

	if !isCorrect {
		return c.JSON(fiber.Map{"result": "Incorrect!"})
	}

	// Calculate points first (before saving submission)
	var pointsEarned int
	if isCorrect {
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
		pointsEarned = calculateDynamicPoints(challenge, solveCount)
		println("=== SubmitFlag Points Calculation ===")
		println("Challenge ID:", challenge.ID)
		println("Base Points:", challenge.Points)
		println("Solve Count:", solveCount)
		println("Points Earned:", pointsEarned)
		println("=====================================")
	}

	// save submission with points earned (chỉ khi hợp lệ)
	sub := models.Submission{
		UserID:       userID,
		TeamID:       user.TeamID, // Include team ID in submission
		ChallengeID:  challenge.ID,
		Flag:         input.Flag,
		IsCorrect:    isCorrect,
		PointsEarned: pointsEarned,
	}
	database.DB.Create(&sub)

	if !existingSolve {
		// Only add to SolvedBy if this is the first solve for user/team
		database.DB.Model(&challenge).Association("SolvedBy").Append(&user)
	}
	response := fiber.Map{"result": "Correct!", "points": pointsEarned}
	if existingSolve {
		response["message"] = "You already solved this challenge before"
	}
	return c.JSON(response)
}
