package middleware

import (
	"log"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/gofiber/fiber/v2"
)

func RoleRequired(allowedRoles ...desc.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("DEBUG header: '%s'", c.Get("Authorization"))
		log.Printf("DEBUG cookie: '%s'", c.Cookies("access_token"))

		// get role from locals
		localsRole := c.Locals("role")
		if localsRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"err": "user role not found",
			})
		}

		// parse to string
		userRole, ok := localsRole.(int32)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"err": "failed to parse user role",
			})
		}

		// check if the role is allowed
		for _, allowedRole := range allowedRoles {
			if userRole == int32(allowedRole) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"err":  "access denied, missing required role",
			"role": userRole,
		})
	}
}
