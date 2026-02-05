package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/utils"
	"context"
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func CreateUser(c *fiber.Ctx) error {
	return GrpcHandler(c, initializers.GrpcClient.CreateUser)
}

func GetAllUsers(c *fiber.Ctx) error {
	return GrpcHandler(c, initializers.GrpcClient.GetAllUsers)
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
	user, err := initializers.GrpcClient.GetUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

func UpdateUser(c *fiber.Ctx) error {
	return GrpcHandler(c, initializers.GrpcClient.UpdateUser)
}

func UpdatePassword(c *fiber.Ctx) error {
	return GrpcHandler(c, initializers.GrpcClient.UpdatePassword)
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
	_, err = initializers.GrpcClient.DeleteUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted user": id})
}

func VerifyPassword(c *fiber.Ctx) error {
	return GrpcHandler(c, initializers.GrpcClient.VerifyPassword)
}

func GrpcHandler[Request any, Response any](
	c *fiber.Ctx,
	call func(context.Context, *Request, ...grpc.CallOption) (*Response, error)) error {
	req := new(Request)
	if err := utils.ParseBodyAndValidate[Request](c, req); err != nil {
		return err
	}

	res, err := call(c.UserContext(), req)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
