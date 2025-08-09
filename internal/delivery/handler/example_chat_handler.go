package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/mappers"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

// ExampleChatHandler shows how to use the new response models and mappers
type ExampleChatHandler struct {
	chatUsecase *usecase.ChatUsecase
}

func NewExampleChatHandler(chatUsecase *usecase.ChatUsecase) *ExampleChatHandler {
	return &ExampleChatHandler{
		chatUsecase: chatUsecase,
	}
}

// GetSessionExample shows how to use mappers to return clean response
func (h *ExampleChatHandler) GetSessionExample(c *fiber.Ctx) error {
	sessionIDStr := c.Params("session_id")
	if sessionIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID is required",
			Error:   "invalid parameter",
		})
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid session ID format",
			Error:   err.Error(),
		})
	}

	// Get entity from usecase (this returns *entities.ChatSession)
	session, err := h.chatUsecase.GetSession(c.Context(), sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get session",
			Error:   err.Error(),
		})
	}

	if session == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session not found",
			Error:   "session does not exist",
		})
	}

	// 🎯 PERBEDAAN UTAMA: Convert entity to clean response using mapper
	response := mappers.ChatSessionToDetailResponse(session)

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Session retrieved successfully",
		Data:    response, // Now returns clean ChatSessionDetailResponse instead of raw entity
	})
}

// GetSessionsExample shows how to use mappers for list response with pagination
func (h *ExampleChatHandler) GetSessionsExample(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get entities from usecase (need to provide required parameters)
	sessions, total, err := h.chatUsecase.GetSessions(c.Context(), page, limit, status, nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to get sessions",
			Error:   err.Error(),
		})
	}

	// 🎯 PERBEDAAN UTAMA: Convert entities to clean response using mapper
	sessionResponses := mappers.ChatSessionsToMinimalResponse(sessions)

	// Create paginated response
	paginatedResponse := mappers.CreatePaginatedResponse(sessionResponses, page, limit, int64(total))

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Sessions retrieved successfully",
		Data:    paginatedResponse, // Now returns clean paginated response
	})
}
