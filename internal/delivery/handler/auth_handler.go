package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} domain.ApiResponse{data=domain.LoginResponse}
// @Failure 400 {object} domain.ApiResponse
// @Failure 401 {object} domain.ApiResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Email and password are required",
			Error:   "validation failed",
		})
	}

	// Get client IP for rate limiting and security
	clientIP := c.IP()
	userAgent := c.Get("User-Agent")

	response, err := h.authUsecase.Login(c.Context(), &req, clientIP, userAgent)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
	}

	// Set secure cookie for refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register godoc
// @Summary Register new user
// @Description Register a new user (admin only)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register request"
// @Success 201 {object} domain.ApiResponse{data=domain.User}
// @Failure 400 {object} domain.ApiResponse
// @Failure 409 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.Email == "" || req.Password == "" || req.Name == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Email, password, name, and role are required",
			Error:   "validation failed",
		})
	}

	if req.Role != "admin" && req.Role != "agent" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Role must be 'admin' or 'agent'",
			Error:   "validation failed",
		})
	}

	user, err := h.authUsecase.Register(c.Context(), &req)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(domain.ApiResponse{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.ApiResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} domain.ApiResponse{data=domain.RefreshTokenResponse}
// @Failure 400 {object} domain.ApiResponse
// @Failure 401 {object} domain.ApiResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Refresh token is required",
			Error:   "validation failed",
		})
	}

	response, err := h.authUsecase.RefreshToken(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "Token refresh failed",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current authenticated user profile
// @Tags Authentication
// @Produce json
// @Success 200 {object} domain.ApiResponse{data=domain.User}
// @Failure 401 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "User not found in context",
			Error:   "authentication required",
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// Logout godoc
// @Summary User logout
// @Description Invalidate user session and tokens
// @Tags Authentication
// @Produce json
// @Success 200 {object} domain.ApiResponse
// @Failure 401 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "User not found in context",
			Error:   "authentication required",
		})
	}

	// Get token from header
	authHeader := c.Get("Authorization")
	token := ""
	if authHeader != "" && len(authHeader) > 7 {
		token = authHeader[7:] // Remove "Bearer " prefix
	}

	// Get refresh token from cookie
	refreshToken := c.Cookies("refresh_token")

	err := h.authUsecase.Logout(c.Context(), user.ID, token, refreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Logout failed",
			Error:   err.Error(),
		})
	}

	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// ValidateSession godoc
// @Summary Validate current session
// @Description Check if current session is valid and return user info
// @Tags Authentication
// @Produce json
// @Success 200 {object} domain.ApiResponse{data=domain.User}
// @Failure 401 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/auth/validate [get]
func (h *AuthHandler) ValidateSession(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid session",
			Error:   "user not found in context",
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Session is valid",
		Data:    user,
	})
}
