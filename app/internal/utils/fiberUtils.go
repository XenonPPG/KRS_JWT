package utils

import (
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

func ParseBodyAndValidate[T any](c *fiber.Ctx, in any) error {
	if err := c.BodyParser(in); err != nil {
		return BadRequest(c)
	}

	if err := Validator.Struct(in); err != nil {
		return BadRequest(c)
	}

	return nil
}

func BadRequest(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "bad request"})
}

func InternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": "internal server error"})
}

func GetTargetId(c *fiber.Ctx) (targetId int64, err error) {
	// admins can alter any user
	// users can alter only themselves
	role := c.Locals("role").(string)
	if role == string(desc.UserRole_ADMIN) {
		targetId, err = strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return
		}
	} else {
		targetId = c.Locals("userID").(int64)
	}

	return
}
