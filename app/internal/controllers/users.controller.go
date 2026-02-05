package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/utils"
	"context"
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	request := &desc.CreateUserRequest{}

	// parse and validate
	err := utils.ParseBodyAndValidate[desc.CreateUserRequest](c, request)
	if err != nil {
		return err
	}

	// hash password
	request.Password, err = utils.HashPassword(request.GetPassword())
	if err != nil {
		return utils.InternalServerError(c)
	}

	// create user
	user, err := initializers.GrpcClient.CreateUser(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "user created", "user": user})
}

func GetAllUsers(c *fiber.Ctx) error {
	request := &desc.GetAllUsersRequest{}

	// parse and validate
	err := utils.ParseBodyAndValidate[desc.GetAllUsersRequest](c, request)
	if err != nil {
		return err
	}

	// get users
	users, err := initializers.GrpcClient.GetAllUsers(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}

func GetUser(c *fiber.Ctx) error {
	request := &desc.GetUserRequest{}

	// parse
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}
	request.Id = int64(id)

	// get user
	user, err := initializers.GrpcClient.GetUser(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

func UpdateUser(c *fiber.Ctx) error {
	request := &desc.UpdateUserRequest{}

	// parse and validate
	err := utils.ParseBodyAndValidate[desc.UpdateUserRequest](c, request)
	if err != nil {
		return err
	}

	// update user
	user, err := initializers.GrpcClient.UpdateUser(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "updated user", "user": user})
}

func UpdatePassword(c *fiber.Ctx) error {
	request := &desc.UpdatePasswordRequest{}

	// parse and validate
	err := utils.ParseBodyAndValidate[desc.UpdatePasswordRequest](c, request)
	if err != nil {
		return err
	}

	// hash password
	request.NewPassword, err = utils.HashPassword(request.GetNewPassword())
	if err != nil {
		return utils.InternalServerError(c)
	}

	// update password
	_, err = initializers.GrpcClient.UpdatePassword(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "password updated"})
}

func DeleteUser(c *fiber.Ctx) error {
	request := &desc.DeleteUserRequest{}

	// parse
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}
	request.Id = int64(id)

	// delete user
	_, err = initializers.GrpcClient.DeleteUser(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "user deleted"})
}

func VerifyPassword(c *fiber.Ctx) error {
	request := &desc.VerifyPasswordRequest{}

	// parse and validate
	err := utils.ParseBodyAndValidate[desc.VerifyPasswordRequest](c, request)
	if err != nil {
		return err
	}

	// verify password
	isValid, err := initializers.GrpcClient.VerifyPassword(context.Background(), request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"valid": isValid.GetValid()})
}
