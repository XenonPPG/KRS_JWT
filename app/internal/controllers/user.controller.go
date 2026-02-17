package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/middleware"
	"JWT/internal/models"
	"JWT/internal/utils"
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	request := desc.CreateUserRequest{}

	// parse request body
	if err := utils.ParseBodyAndValidate[desc.CreateUserRequest](c, &request); err != nil {
		return utils.BadRequest(c)
	}

	// hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.InternalServerError(c)
	}

	request.Password = hashedPassword
	user, err := initializers.GrpcUserService.CreateUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"user": user})
}

func GetAllUsers(c *fiber.Ctx) error {
	request := models.GetAllUsersRequest{}

	// parse query
	if err := c.QueryParser(&request); err != nil {
		return utils.BadRequest(c)
	}

	users, err := initializers.GrpcUserService.GetAllUsers(c.UserContext(), &desc.GetAllUsersRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}

func GetUser(c *fiber.Ctx) error {
	request := desc.GetUserRequest{}

	// parse
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}
	request.Id = int64(id)

	// get user
	user, err := initializers.GrpcUserService.GetUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

func UpdateUser(c *fiber.Ctx) error {
	return utils.GrpcHandler(c, initializers.GrpcUserService.UpdateUser)
}

func UpdatePassword(c *fiber.Ctx) error {
	return utils.GrpcHandler(c, initializers.GrpcUserService.UpdatePassword)
}

func DeleteUser(c *fiber.Ctx) error {
	request := desc.DeleteUserRequest{}

	// parse
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c)
	}
	request.Id = int64(id)

	// get user
	_, err = initializers.GrpcUserService.DeleteUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted user": id})
}

func Login(c *fiber.Ctx) error {
	request := desc.LoginRequest{}

	// parse request
	if err := c.BodyParser(&request); err != nil {
		return utils.BadRequest(c)
	}

	// get response
	response, err := initializers.GrpcUserService.Login(c.Context(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	// make token
	access, refresh, err := middleware.GenerateTokenPair(response.Id, response.Role)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access": access, "refresh": refresh})
}
