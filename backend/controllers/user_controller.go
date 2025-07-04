package controllers

import (
	"ctfme/database"
	"ctfme/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
)

func RegisterUser(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username	string `json:"username"`
		Email		string `json:"email"`
		Password	string `json:"password"`
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
		Username : req.Username,
		Email : req.Email,
		Password : string(hashedPassword),
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
		"user_id": user.ID,
		"username": user.Username,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create token"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
}

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
		Email		string 	`json:"email"`
		IsAdmin		bool 	`json:"is_admin"`
		TeamID		*uint 	`json:"team_id"`
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user.Email = req.Email
	user.IsAdmin = req.IsAdmin
	user.TeamID = req.TeamID

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update user"})
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
	if err := database.DB.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted!"})
}