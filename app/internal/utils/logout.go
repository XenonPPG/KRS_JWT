package utils

import (
	"JWT/internal/initializers"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		_ = initializers.TokenService.DeleteRefreshToken(c.Context(), refreshToken)
	}

	// delete cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 999), // negative time deletes cookie
		HTTPOnly: false,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 999),
		HTTPOnly: true,
		Path:     "/api/auth/refresh",
	})
}
