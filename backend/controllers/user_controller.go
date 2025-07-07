package controllers

import (
	"ctfme/config"
	"ctfme/database"
	"ctfme/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User already exists!"})
	}

	return c.JSON(fiber.Map{"message": "Register success!"})
}

func LoginUser(c *fiber.Ctx) error {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User not found!"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid password!"})
	}

	// create jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create token"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
}

// func LogoutUser(c *fiber.Ctx) error {
// 	userID := c.Locals("user_id").(uint)
// }

func AdminGetUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if !user.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}

	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch users"})
	}

	return c.JSON(users)
}

func AdminGetUserDetail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var admin models.User
	if err := database.DB.First(&admin, userID).Error; err != nil || !admin.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only!"})
	}
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found!"})
	}

	return c.JSON(user)
}

func AdminUpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var admin models.User
	if err := database.DB.First(&admin, userID).Error; err != nil || !admin.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only!"})
	}

	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found!"})
	}

	type UpdateUserRequest struct {
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin"`
		TeamID  *uint  `json:"team_id"`
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Store old team ID before updating
	var oldTeamID *uint
	if user.TeamID != nil {
		oldTeamID = user.TeamID
	}

	user.Email = req.Email
	user.IsAdmin = req.IsAdmin
	user.TeamID = req.TeamID

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update user"})
	}

	// Check if old team is now empty and delete it if so
	if oldTeamID != nil && (req.TeamID == nil || *oldTeamID != *req.TeamID) {
		if err := CheckAndDeleteEmptyTeam(*oldTeamID); err != nil {
			// Log error but don't fail the request
			println("Error checking empty team:", err.Error())
		}
	}

	return c.JSON(fiber.Map{"message": "User updated!", "user": user})
}

func AdminDeleteUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var admin models.User
	if err := database.DB.First(&admin, userID).Error; err != nil || !admin.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only!"})
	}

	id := c.Params("id")

	// Use transaction to ensure data consistency
	return database.DB.Transaction(func(tx *gorm.DB) error {
		println("=== AdminDeleteUser Debug ===")
		println("Deleting user ID:", id)

		// First, find the user
		var user models.User
		if err := tx.First(&user, id).Error; err != nil {
			println("User not found:", err.Error())
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		println("Found user:", user.Username)

		// Delete all submissions by this user first
		println("Deleting submissions...")
		if err := tx.Exec("DELETE FROM submissions WHERE user_id = ?", id).Error; err != nil {
			println("Error deleting submissions:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not delete user submissions"})
		}
		println("Submissions deleted successfully")

		// Clear the many-to-many relationship with challenges (challenge_solvers table)
		println("Clearing challenge solvers...")
		if err := tx.Exec("DELETE FROM challenge_solvers WHERE user_id = ?", id).Error; err != nil {
			println("Error clearing challenge solvers:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not clear user challenge solvers"})
		}
		println("Challenge solvers cleared successfully")

		// Store team ID before removing user from team
		var teamID *uint
		if user.TeamID != nil {
			teamID = user.TeamID
		}

		// Remove user from team (set team_id to NULL)
		println("Removing user from team...")
		if err := tx.Exec("UPDATE users SET team_id = NULL WHERE id = ?", id).Error; err != nil {
			println("Error removing user from team:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not remove user from team"})
		}
		println("User removed from team successfully")

		// Check if team is now empty and delete it if so
		if teamID != nil {
			println("Checking if team is empty...")
			var userCount int64
			if err := tx.Model(&models.User{}).Where("team_id = ?", teamID).Count(&userCount).Error; err != nil {
				println("Error counting team users:", err.Error())
			} else {
				println("Team user count:", userCount)
				if userCount == 0 {
					println("Deleting empty team...")
					if err := tx.Exec("DELETE FROM teams WHERE id = ?", teamID).Error; err != nil {
						println("Error deleting empty team:", err.Error())
					} else {
						println("Empty team deleted successfully")
					}
				}
			}
		}

		// Finally delete the user
		println("Deleting user...")
		if err := tx.Exec("DELETE FROM users WHERE id = ?", id).Error; err != nil {
			println("Error deleting user:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not delete user"})
		}
		println("User deleted successfully")
		println("=============================")

		return c.JSON(fiber.Map{"message": "User deleted successfully!"})
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	var solved []models.Submission
	database.DB.Where("user_id = ? AND is_correct = true", userID).Find(&solved)
	return c.JSON(fiber.Map{"user": user, "solved": solved})
}

func GetScoreboard(c *fiber.Ctx) error {
	println("=== GetScoreboard Debug ===")
	println("TeamMode:", config.TeamMode)
	println("DynamicScoreEnabled:", config.DynamicScoreEnabled)
	println("DynamicScoreDecay:", config.DynamicScoreDecay)
	println("DynamicScoreMin:", config.DynamicScoreMin)

	if config.TeamMode {
		type TeamScore struct {
			TeamName string
			Points   int
		}
		var scores []TeamScore

		// Get all teams with their solved challenges
		var teams []models.Team
		database.DB.Preload("Users.Submissions", "is_correct = true").Find(&teams)

		for _, team := range teams {
			totalPoints := 0
			println("Processing team:", team.Name)

			// Calculate points from submissions (using points earned at submission time)
			challengePoints := make(map[uint]int) // challenge_id -> points earned
			for _, user := range team.Users {
				for _, submission := range user.Submissions {
					if submission.IsCorrect && submission.PointsEarned > 0 {
						// Use the highest points earned for this challenge by this team
						if submission.PointsEarned > challengePoints[submission.ChallengeID] {
							challengePoints[submission.ChallengeID] = submission.PointsEarned
						}
					}
				}
			}
			println("Team solved challenges count:", len(challengePoints))

			// Sum up points earned for each challenge
			for challengeID, pointsEarned := range challengePoints {
				println("Challenge ID:", challengeID, "Points Earned:", pointsEarned)
				totalPoints += pointsEarned
			}
			println("Team total points:", totalPoints)

			scores = append(scores, TeamScore{
				TeamName: team.Name,
				Points:   totalPoints,
			})
		}

		// Sort by points descending
		for i := 0; i < len(scores)-1; i++ {
			for j := i + 1; j < len(scores); j++ {
				if scores[i].Points < scores[j].Points {
					scores[i], scores[j] = scores[j], scores[i]
				}
			}
		}
		println("=============================")

		return c.JSON(scores)
	} else {
		type UserScore struct {
			Username string
			Points   int
		}
		var scores []UserScore

		// Get all users with their solved challenges
		var users []models.User
		database.DB.Preload("Submissions", "is_correct = true").Where("is_admin = false").Find(&users)

		for _, user := range users {
			totalPoints := 0
			challengePoints := make(map[uint]int) // challenge_id -> points earned

			// Calculate points from submissions (using points earned at submission time)
			for _, submission := range user.Submissions {
				if submission.IsCorrect && submission.PointsEarned > 0 {
					// Use the highest points earned for this challenge by this user
					if submission.PointsEarned > challengePoints[submission.ChallengeID] {
						challengePoints[submission.ChallengeID] = submission.PointsEarned
					}
				}
			}

			// Sum up points earned for each challenge
			for _, pointsEarned := range challengePoints {
				totalPoints += pointsEarned
			}

			scores = append(scores, UserScore{
				Username: user.Username,
				Points:   totalPoints,
			})
		}

		// Sort by points descending
		for i := 0; i < len(scores)-1; i++ {
			for j := i + 1; j < len(scores); j++ {
				if scores[i].Points < scores[j].Points {
					scores[i], scores[j] = scores[j], scores[i]
				}
			}
		}

		return c.JSON(scores)
	}
}
