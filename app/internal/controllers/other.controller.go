package controllers

import "github.com/gofiber/fiber/v2"

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags health
// @Produce plain
// @Success 200 {string} string "Healthy :)"
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).SendString("Healthy :)")
}
