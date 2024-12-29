package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strings"
)

var jwtSecret = "your_secret_key"

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the Authorization header
	authHeader := c.Get("Authorization")

	// Check if the header is missing or doesn't start with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
	}

	// Extract the token part
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	//	tokenString := c.Get("Authorization")[7:] // Skip "Bearer "
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	c.Locals("user", token)

	// Extract user type from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}
	userType := claims["UserType"]
	fmt.Println("Extracted userType:", userType)
	c.Locals("userType", userType)
	return c.Next()
}
