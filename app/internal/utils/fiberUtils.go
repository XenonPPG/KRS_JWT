package utils

import (
	"fmt"
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

func InternalServerError(c *fiber.Ctx, message error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"err": "internal server error",
		"msg": message.Error(),
	})
}

func GetTargetId(c *fiber.Ctx) (targetId int64, err error) {
	role, ok := GetLocalAndParse[int32](c, "role")
	if !ok {
		return 0, fmt.Errorf("role not found in context")
	}

	if role == int32(desc.UserRole_ADMIN) {
		targetId, err = strconv.ParseInt(c.Params("id"), 10, 64)
	} else {
		targetId, ok = GetLocalAndParse[int64](c, "user_id")
		if !ok {
			err = fmt.Errorf("user_id not found in context")
		}
	}

	return
}

func GetLocalAndParse[T any](c *fiber.Ctx, key string) (out T, ok bool) {
	item := c.Locals(key)
	if item == nil {
		return out, false
	}

	out, ok = item.(T)
	return
}
