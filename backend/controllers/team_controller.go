package controllers

import (
	"crypto/rand"
	"ctfme/database"
	"ctfme/models"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func JoinTeam(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var requests struct {
		InviteCode string `json:"invite_code"`
	}

	if err := c.BodyParser(&requests); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var team models.Team
	if err := database.DB.Where("invite_code = ?", requests.InviteCode).First(&team).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Team not found"})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if user is already in this team
	if user.TeamID != nil && *user.TeamID == team.ID {
		return c.Status(400).JSON(fiber.Map{"error": "You are already in this team"})
	}

	// Store old team ID before joining new team
	var oldTeamID *uint
	if user.TeamID != nil {
		oldTeamID = user.TeamID
	}

	user.TeamID = &team.ID
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not join team"})
	}

	// Check if old team is now empty and delete it if so
	if oldTeamID != nil {
		if err := CheckAndDeleteEmptyTeam(*oldTeamID); err != nil {
			// Log error but don't fail the request
			println("Error checking empty team:", err.Error())
		}
	}

	return c.JSON(fiber.Map{"message": "Joined team successfully!", "team": team})
}

func LeaveTeam(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Store team ID before removing user from team
	var teamID *uint
	if user.TeamID != nil {
		teamID = user.TeamID
	}

	user.TeamID = nil
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not leave team"})
	}

	// Check if team is now empty and delete it if so
	if teamID != nil {
		if err := CheckAndDeleteEmptyTeam(*teamID); err != nil {
			// Log error but don't fail the request
			println("Error checking empty team:", err.Error())
		}
	}

	return c.JSON(fiber.Map{"message": "Left team successfully!"})
}

func CreateTeam(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	// Use transaction to ensure data consistency
	return database.DB.Transaction(func(tx *gorm.DB) error {
		println("=== CreateTeam Debug ===")
		println("UserID:", userID)
		println("Team Name:", req.Name)

		// First, check if user exists
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			println("User not found:", err.Error())
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		println("Found user:", user.Username)
		println("User current team ID:", user.TeamID)

		// Check if user is already in a team
		if user.TeamID != nil {
			println("User already in team:", *user.TeamID)
			return c.Status(400).JSON(fiber.Map{"error": "You are already in a team. Please leave your current team first."})
		}
		println("User is not in any team, proceeding...")

		// Check if team name already exists
		var exists models.Team
		if err := tx.Where("name = ?", req.Name).First(&exists).Error; err == nil {
			println("Team name already exists:", req.Name)
			return c.Status(400).JSON(fiber.Map{"error": "Team name already exists"})
		}
		println("Team name is available")

		// Generate invite code
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Could not generate invite code"})
		}
		inviteCode := hex.EncodeToString(b)

		// Create team
		team := models.Team{
			Name:       req.Name,
			InviteCode: inviteCode,
		}
		if err := tx.Create(&team).Error; err != nil {
			println("Error creating team:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not create team"})
		}
		println("Team created with ID:", team.ID)

		// Add user to team
		user.TeamID = &team.ID
		if err := tx.Save(&user).Error; err != nil {
			println("Error adding user to team:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not add user to team"})
		}
		println("User added to team successfully")

		// Preload users for response
		if err := tx.Preload("Users").First(&team, team.ID).Error; err != nil {
			println("Error loading team details:", err.Error())
			return c.Status(500).JSON(fiber.Map{"error": "Could not load team details"})
		}
		println("Team details loaded, user count:", len(team.Users))
		println("=============================")

		return c.JSON(fiber.Map{"message": "Team created successfully!", "team": team})
	})
}

// API for admin to get all teams
func AdminGetTeams(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var admin models.User
	if err := database.DB.First(&admin, userID).Error; err != nil || !admin.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	var teams []models.Team
	if err := database.DB.Preload("Users").Find(&teams).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch teams"})
	}
	return c.JSON(teams)
}

// API for admin to get team detail
func AdminGetTeamDetail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var admin models.User
	if err := database.DB.First(&admin, userID).Error; err != nil || !admin.IsAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
	}
	teamID := c.Params("id")
	var team models.Team
	if err := database.DB.Preload("Users").First(&team, teamID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Team not found"})
	}
	// ensure Users is always an array (not null)
	if team.Users == nil {
		team.Users = []models.User{}
	}
	return c.JSON(team)
}

// Helper function to check and delete empty teams
func CheckAndDeleteEmptyTeam(teamID uint) error {
	var userCount int64
	if err := database.DB.Model(&models.User{}).Where("team_id = ?", teamID).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
		// Team is empty, delete it
		if err := database.DB.Unscoped().Delete(&models.Team{}, teamID).Error; err != nil {
			return err
		}
		println("Empty team deleted:", teamID)
	}

	return nil
}
