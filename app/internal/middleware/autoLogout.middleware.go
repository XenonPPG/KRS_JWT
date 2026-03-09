package middleware

import (
	"JWT/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AutoLogout(c *fiber.Ctx) error {
	utils.Logout(c)

	return c.Next()
}
