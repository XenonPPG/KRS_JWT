package main

import (
	"JWT/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("/api")
	api.Route("/user", func(router fiber.Router) {
		router.Post("/create", controllers.CreateUser)
		router.Get("/", controllers.GetAllUsers)
		router.Get("/:id", controllers.GetUser)
		router.Put("/", controllers.UpdateUser)
		router.Put("/password", controllers.UpdatePassword)
		router.Delete("/:id", controllers.DeleteUser)
		router.Get("/verify", controllers.VerifyPassword)
	})
}
