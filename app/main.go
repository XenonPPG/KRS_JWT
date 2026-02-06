package main

import (
	"JWT/internal/controllers"
	"JWT/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("/api")
	api.Route("/user", func(router fiber.Router) {
		router.Post("/", controllers.CreateUser)
		router.Get("/:id", controllers.GetUser)
		router.Get("/verify", controllers.VerifyPassword)

		protected := router.Group("/", middleware.JWTProtected)
		protected.Delete("/:id", controllers.DeleteUser)
		protected.Get("/", controllers.GetAllUsers)
		protected.Put("/", controllers.UpdateUser)
		protected.Put("/password", controllers.UpdatePassword)
	})
}
