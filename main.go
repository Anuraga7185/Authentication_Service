package main

import (
	"github.com/Anuraga7185/Libraries/authservice/database"
	"github.com/Anuraga7185/Libraries/authservice/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// Initialize the database
	database.ConnectDB()

	// Create a Fiber app
	app := fiber.New()

	// Register routes
	routes.AuthRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
