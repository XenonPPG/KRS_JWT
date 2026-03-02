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

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with hashed password
// @Tags users
// @Accept json
// @Produce json
// @Param request body desc.CreateUserRequest true "User creation request"
// @Success 201 {object} map[string]interface{} "user created successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	request := desc.CreateUserRequest{}

	// parse request body
	if err := utils.ParseBodyAndValidate[desc.CreateUserRequest](c, &request); err != nil {
		return utils.BadRequest(c)
	}

	// hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	request.Password = hashedPassword
	user, err := initializers.GrpcUserService.CreateUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"user": user})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} map[string]interface{} "users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /users [get]
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
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves a specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "user retrieved successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /users/{id} [get]
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
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Updates an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body desc.UpdateUserRequest true "User update request"
// @Success 200 {object} map[string]interface{} "user updated successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /users [put]
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
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"updated user": response})
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Updates the password for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "password updated successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /users/password [put]
func UpdatePassword(c *fiber.Ctx) error {
	return utils.GrpcHandler(c, initializers.GrpcUserService.UpdatePassword)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "user deleted successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	request := desc.DeleteUserRequest{}

	targetId, err := utils.GetTargetId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Bad request",
			"msg": err.Error(),
		})
	}
	request.Id = targetId

	// get user
	_, err = initializers.GrpcUserService.DeleteUser(c.UserContext(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted user": targetId})
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns JWT access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body desc.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "tokens issued successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	request := desc.LoginRequest{}

	// parse request
	if err := c.BodyParser(&request); err != nil {
		return utils.BadRequest(c)
	}

	// get response
	response, err := initializers.GrpcUserService.Login(c.Context(), &request)
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	// make token
	access, refresh, err := middleware.IssueTokenPair(c, &models.UserInfo{
		Id:   response.Id,
		Role: *response.Role,
	})
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access": access, "refresh": refresh})
}

// Logout godoc
// @Summary User logout
// @Description Logs out the user by invalidating refresh token and clearing cookies
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{} "logged out successfully"
// @Router /auth/logout [post]
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

// RefreshTokens godoc
// @Summary Refresh access tokens
// @Description Rotates refresh token and issues new access and refresh tokens
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{} "tokens rotated successfully"
// @Failure 401 {object} map[string]interface{} "unauthorized - invalid or missing refresh token"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /auth/refresh [post]
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
		return utils.InternalServerError(c, err)
	}

	// issue new tokens
	_, _, err = middleware.IssueTokenPair(c, userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": "could not generate new tokens"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "tokens rotated"})
}
