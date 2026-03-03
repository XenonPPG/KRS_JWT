package controllers

import "github.com/gofiber/fiber/v2"

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "status"
// @Router /api/health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "Healthy :)",
	})
}
