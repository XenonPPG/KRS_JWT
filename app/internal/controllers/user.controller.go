package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/models"
	"JWT/internal/utils"
	"strconv"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/gofiber/fiber/v2"
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
// @Router /api/user [post]
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
// @Param ascending_order query bool false "Sort in ascending order"
// @Success 200 {object} map[string]interface{} "users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/user [get]
func GetAllUsers(c *fiber.Ctx) error {
	request := models.GetAllItemsRequest{}

	// parse query
	if err := c.QueryParser(&request); err != nil {
		return utils.BadRequest(c)
	}

	users, err := initializers.GrpcUserService.GetAllUsers(c.UserContext(), &desc.GetAllUsersRequest{
		Limit:          request.Limit,
		Offset:         request.Offset,
		AscendingOrder: request.AscendingOrder,
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
// @Router /api/user/{id} [get]
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
// @Router /api/user/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	var request desc.UpdateUserRequest

	targetId, err := utils.GetTargetId(c)
	if err != nil {
		return utils.BadRequest(c)
	}

	// parse request body
	if err = utils.ParseBodyAndValidate[desc.UpdateUserRequest](c, &request); err != nil {
		return utils.BadRequest(c)
	}
	request.Id = targetId

	// get user
	user, err := initializers.GrpcUserService.GetUser(c.UserContext(), &desc.GetUserRequest{Id: targetId})
	if err != nil {
		return utils.InternalServerError(c, err)
	}

	// admins cannot be downgraded
	if user.GetRole() == desc.UserRole_ADMIN && int32(request.GetRole()) < int32(user.GetRole()) {
		return utils.BadRequest(c)
	}

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
// @Param request body desc.UpdatePasswordRequest true "Password update request"
// @Success 200 {object} map[string]interface{} "password updated successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/user/password [put]
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
// @Router /api/user/{id} [delete]
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
