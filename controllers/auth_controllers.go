package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Anuraga7185/Libraries/authservice/database"
	"github.com/Anuraga7185/Libraries/authservice/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

/*import (
"auth-service/database"
"auth-service/models"
"context"
"crypto/sha256"
"encoding/hex"
"time"

"github.com/gofiber/fiber/v2"
)*/

func SignUp(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Connect to the users collection
	collection := database.DB.Collection("users")

	// Check if a user with the same email already exists
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		// User with this email already exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	if err != mongo.ErrNoDocuments {
		// Unexpected error during the database query
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking for existing user",
		})
	}

	// Hash the password
	hash := sha256.New()
	hash.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserType:  req.UserType,
	}

	// Save to MongoDB
	//collection := database.DB.Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}
func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Connect to the MongoDB collection
	collection := database.DB.Collection("users")

	// Find user by email
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Hash the password for comparison
	hash := sha256.New()
	hash.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	if user.Password != hashedPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT
	claims := jwt.MapClaims{
		"email":    user.Email,
		"UserType": user.UserType, // Store user type in the JWT
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}
func Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	return c.JSON(fiber.Map{"message": "Welcome!", "email": email})
}

func GetUsers(c *fiber.Ctx) error { // Select the database and collection
	// Get userType from the context (set by Auth middleware)
	userType := c.Locals("userType").(string)
	collection := database.DB.Collection("users")

	// Check if the user is an Admin
	if userType != "ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied. Only Admins can view all users.",
		})
	}
	// If user is an Admin, fetch all users from the collection (without filtering by user_type)
	var users []models.User
	filter := bson.M{}
	cursor, err := collection.Find(context.Background(), filter) // Empty filter to get all users
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching all users.",
		})
	}
	defer cursor.Close(context.Background())

	// Decode each user document into the `users` slice
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error decoding user data")
		}
		users = append(users, user)
	}

	// Check for cursor iteration error
	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error iterating cursor")
	}

	// Return all users as JSON
	return c.JSON(users)
}
