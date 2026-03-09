package controllers

import (
	"JWT/internal/initializers"
	"JWT/internal/middleware"
	"JWT/internal/models"
	"JWT/internal/utils"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns JWT access, refresh tokens, user info
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body desc.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "tokens issued successfully"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/auth/login [post]
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access":  access,
		"refresh": refresh,
		"user":    response,
	})
}

// LogoutHandler godoc
// @Summary User logout
// @Description Logs out the user by invalidating refresh token and clearing cookies
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]interface{} "logged out successfully"
// @Router /api/auth/logout [post]
func LogoutHandler(c *fiber.Ctx) error {
	utils.Logout(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": "logged out"})
}

// RefreshTokens godoc
// @Summary Refresh access tokens
// @Description Rotates refresh token and issues new access and refresh tokens
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]interface{} "tokens rotated successfully"
// @Failure 401 {object} map[string]interface{} "unauthorized - invalid or missing refresh token"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/auth/refresh [post]
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
