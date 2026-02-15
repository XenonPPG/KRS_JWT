package utils

import (
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
