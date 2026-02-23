package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/middleware"
	"JWT/internal/models"
	"JWT/internal/utils"
	"strconv"
	"time"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	request := models.GetAllItemsRequest{}

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
	var request desc.UpdateUserRequest

	targetId, err := utils.GetTargetId(c)
	if err != nil {
		return utils.BadRequest(c)
	}

	// parse request body
	if err := utils.ParseBodyAndValidate[desc.UpdateUserRequest](c, &request); err != nil {
		return utils.BadRequest(c)
	}
	request.Id = targetId

	// update user
	response, err := initializers.GrpcUserService.UpdateUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"updated user": response})
}

func UpdatePassword(c *fiber.Ctx) error {
	return utils.GrpcHandler(c, initializers.GrpcUserService.UpdatePassword)
}

func DeleteUser(c *fiber.Ctx) error {
	request := desc.DeleteUserRequest{}

	targetId, err := utils.GetTargetId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "Bad request"})
	}
	request.Id = targetId

	// get user
	_, err = initializers.GrpcUserService.DeleteUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted user": targetId})
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
	access, refresh, err := middleware.IssueTokenPair(c, &models.UserInfo{
		Id:   response.Id,
		Role: *response.Role,
	})
	if err != nil {
		return utils.InternalServerError(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access": access, "refresh": refresh})
}

func Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		_ = initializers.TokenService.DeleteRefreshToken(c.Context(), refreshToken)
	}

	// delete cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 999), // Прошедшее время удаляет куку
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 999),
		HTTPOnly: true,
		Path:     "/api/auth/refresh",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "logged out"})
}

func RefreshTokens(c *fiber.Ctx) error {
	// get refresh token from cookie
	oldRefresh := c.Cookies("refresh_token")
	if oldRefresh == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": "refresh token missing"})
	}

	// check jwt validity
	_, err := jwt.Parse(oldRefresh, func(token *jwt.Token) (interface{}, error) {
		return initializers.RefreshSecret, nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": "invalid or expired refresh token"})
	}

	// check refresh token in storage
	userInfo, err := initializers.TokenService.ValidateRefreshToken(c.Context(), oldRefresh)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": "session not found or expired in storage"})
	}

	// delete old refresh token
	err = initializers.TokenService.DeleteRefreshToken(c.Context(), oldRefresh)
	if err != nil {
		return utils.InternalServerError(c)
	}

	// issue new tokens
	_, _, err = middleware.IssueTokenPair(c, userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": "could not generate new tokens"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "tokens rotated"})
}
