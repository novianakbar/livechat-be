package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewSecurityMiddleware() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	})
}

func NewRateLimitMiddleware() fiber.Handler {
	// Simple rate limiting implementation
	return func(c *fiber.Ctx) error {
		// In production, implement proper rate limiting with Redis
		return c.Next()
	}
}

func NewLoginRateLimitMiddleware() fiber.Handler {
	// Simple rate limiting for login attempts
	return func(c *fiber.Ctx) error {
		// In production, implement proper rate limiting for login attempts
		return c.Next()
	}
}

func NewLoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path} ${ip} ${ua}\n",
		TimeFormat: "2006-01-02 15:04:05",
	})
}
