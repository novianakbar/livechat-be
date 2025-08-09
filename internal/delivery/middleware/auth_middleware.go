package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type AuthMiddleware struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthMiddleware(authUsecase *usecase.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authUsecase: authUsecase,
	}
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
				Success: false,
				Message: "Authorization header required",
				Error:   "missing authorization header",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
				Success: false,
				Message: "Invalid authorization header format",
				Error:   "authorization header must be 'Bearer <token>'",
			})
		}

		token := tokenParts[1]
		user, err := m.authUsecase.ValidateToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
				Success: false,
				Message: "Invalid or expired token",
				Error:   err.Error(),
			})
		}

		// Store user in context
		c.Locals("user", user)
		return c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*domain.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
				Success: false,
				Message: "Authentication required",
				Error:   "user not found in context",
			})
		}

		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(domain.ApiResponse{
			Success: false,
			Message: "Insufficient permissions",
			Error:   "user role not authorized for this action",
		})
	}
}

func (m *AuthMiddleware) RequireAgent() fiber.Handler {
	return m.RequireRole("agent", "admin")
}

func (m *AuthMiddleware) RequireAdmin() fiber.Handler {
	return m.RequireRole("admin")
}

func GetUserFromContext(c *fiber.Ctx) *domain.User {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserIDFromContext(c *fiber.Ctx) *string {
	user := GetUserFromContext(c)
	if user == nil {
		return nil
	}
	return &user.ID
}
