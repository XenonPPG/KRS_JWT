package middleware

import (
	"JWT/internal/models"
	"context"
	"fmt"
	"strings"
	"time"

	"JWT/internal/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	Role   int32 `json:"role"`
	jwt.RegisteredClaims
}

func IssueTokenPair(c *fiber.Ctx, userInfo *models.UserInfo) (access string, refresh string, err error) {
	// access Token
	accessClaims := &Claims{
		UserID: userInfo.Id,
		Role:   int32(userInfo.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(initializers.AccessSecret)
	if err != nil {
		return "", "", err
	}

	// refresh Token
	refreshClaims := &jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(initializers.RefreshSecret)
	if err != nil {
		return "", "", err
	}

	// save to redis
	err = initializers.TokenService.StoreRefreshToken(context.Background(), refresh, userInfo, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	// save to cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		Path:     "api/auth",
	})

	return access, refresh, nil
}

// JWTProtected checks only access token
func JWTProtected(c *fiber.Ctx) error {
	tokenString := c.Cookies("access_token")

	if tokenString == "" {
		tokenString = strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": "missing token"})
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return initializers.AccessSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": "access token expired or invalid"})
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"err": "invalid claims"})
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)

	return c.Next()
}
