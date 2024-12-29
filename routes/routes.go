package routes

import (
	"github.com/Anuraga7185/Libraries/authservice/controllers"
	middlewares "github.com/Anuraga7185/Libraries/authservice/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/signup", controllers.SignUp)
	auth.Post("/login", controllers.Login)
	auth.Get("/users", middlewares.AuthMiddleware, controllers.GetUsers)
	auth.Get("/profile", middlewares.AuthMiddleware, controllers.Profile)
}
