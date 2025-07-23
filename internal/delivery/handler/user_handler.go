package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with pagination and filters
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param role query string false "User role filter"
// @Param department_id query string false "Department ID filter"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.User}
// @Failure 500 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	role := c.Query("role")
	departmentIDStr := c.Query("department_id")

	var departmentID *uuid.UUID
	if departmentIDStr != "" {
		id, err := uuid.Parse(departmentIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
				Success: false,
				Message: "Invalid department ID format",
				Error:   err.Error(),
			})
		}
		departmentID = &id
	}

	users, total, err := h.userUsecase.GetUsers(c.Context(), page, limit, role, departmentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get users",
			Error:   err.Error(),
		})
	}

	totalPages := (total + limit - 1) / limit
	pagination := domain.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return c.JSON(domain.PaginatedResponse{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

// GetAgents godoc
// @Summary Get all agents
// @Description Get all users with agent role
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} domain.ApiResponse{data=[]domain.User}
// @Failure 500 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/users/agents [get]
func (h *UserHandler) GetAgents(c *fiber.Ctx) error {
	agents, err := h.userUsecase.GetAgents(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get agents",
			Error:   err.Error(),
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Agents retrieved successfully",
		Data:    agents,
	})
}

// GetUser godoc
// @Summary Get single user
// @Description Get a single user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.ApiResponse{data=domain.User}
// @Failure 400 {object} domain.ApiResponse
// @Failure 404 {object} domain.ApiResponse
// @Security BearerAuth
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "User ID is required",
			Error:   "invalid parameter",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
	}

	user, err := h.userUsecase.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get user",
			Error:   err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ApiResponse{
			Success: false,
			Message: "User not found",
			Error:   "user does not exist",
		})
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}
