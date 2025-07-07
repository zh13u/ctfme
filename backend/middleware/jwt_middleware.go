package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		println("=== JWT Middleware Debug ===")
		println("Request path:", c.Path())

		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			println("Missing or invalid Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}

		tokenStr := authHeader[7:]
		println("Token length:", len(tokenStr))

		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			println("Token parse error:", err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}
		if !token.Valid {
			println("Token is invalid")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}
		println("Token is valid")
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			println("Invalid token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			println("Invalid user_id in token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user_id in token"})
		}
		println("User ID from token:", uint(userID))
		c.Locals("user_id", uint(userID))
		println("JWT Middleware completed successfully")
		println("=============================")
		return c.Next()
	}
}
