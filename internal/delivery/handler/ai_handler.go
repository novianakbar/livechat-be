package handler

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/delivery/middleware"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/service"
	"github.com/novianakbar/livechat-be/internal/usecase"
)

type AIHandler struct {
	chatUsecase  *usecase.ChatUsecase
	kafkaService *service.KafkaService
}

func NewAIHandler(chatUsecase *usecase.ChatUsecase, kafkaService *service.KafkaService) *AIHandler {
	return &AIHandler{
		chatUsecase:  chatUsecase,
		kafkaService: kafkaService,
	}
}

// ReceiveAIResponse godoc
// @Summary Receive response from AI system
// @Description Endpoint for AI system to send responses back to livechat system
// @Tags ai
// @Accept json
// @Produce json
// @Param request body domain.SendMessageRequest true "AI response message"
// @Success 200 {object} domain.ApiResponse{data=domain.SendMessageResponse}
// @Failure 400 {object} domain.ApiResponse
// @Failure 401 {object} domain.ApiResponse
// @Security ApiKeyAuth
// @Router /api/ai/response [post]
func (h *AIHandler) ReceiveAIResponse(c *fiber.Ctx) error {
	var req domain.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.SessionID == uuid.Nil || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID and message are required",
			Error:   "validation failed",
		})
	}

	// Get user from context (AI user that was authenticated)
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "Authentication required",
			Error:   "user not found in context",
		})
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
	}

	// Send message with AI as sender
	response, err := h.chatUsecase.SendMessage(c.Context(), &req, &userUUID, "ai")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to send AI response",
			Error:   err.Error(),
		})
	}

	// Kirim typing stop sebelum mempublikasikan pesan
	// Ini untuk memastikan indikator typing berhenti saat pesan AI diterima
	typingMessage := domain.TypingMessage{
		Type:      "typing_indicator",
		SessionID: req.SessionID,
		UserID:    user.ID,
		UserType:  "ai",
		IsTyping:  false, // set false = typing stop
		Timestamp: time.Now(),
	}

	// Kirim typing stop ke Kafka
	if err := h.kafkaService.PublishTypingIndicator(c.Context(), typingMessage); err != nil {
		// Hanya log error, jangan hentikan proses
		log.Printf("Failed to send typing stop indicator for AI: %v", err)
	} else {
		log.Printf("Typing stop indicator sent for AI on session %s", req.SessionID)
	}

	// Get the message from the database to broadcast via Kafka, similar to ChatHandler.SendMessage
	message, err := h.chatUsecase.GetMessageByID(c.Context(), response.MessageID)

	// Publish message to Kafka to deliver to the user
	if h.kafkaService != nil && message != nil {
		// Parse IDs for Kafka message
		messageUUID, _ := uuid.Parse(message.ID)
		sessionUUID, _ := uuid.Parse(message.SessionID)

		kafkaMessage := struct {
			ID          uuid.UUID  `json:"id"`
			SessionID   uuid.UUID  `json:"session_id"`
			SenderID    *uuid.UUID `json:"sender_id"`
			SenderType  string     `json:"sender_type"`
			Message     string     `json:"message"`
			MessageType string     `json:"message_type"`
			Attachments []string   `json:"attachments"`
			ReadAt      *time.Time `json:"read_at"`
			CreatedAt   time.Time  `json:"created_at"`
			UpdatedAt   time.Time  `json:"updated_at"`
		}{
			ID:        messageUUID,
			SessionID: sessionUUID,
			SenderID: func() *uuid.UUID {
				if message.SenderID.Valid {
					if parsed, err := uuid.Parse(message.SenderID.String); err == nil {
						return &parsed
					}
				}
				return nil
			}(),
			SenderType:  message.SenderType,
			Message:     message.Message,
			MessageType: message.MessageType,
			Attachments: message.Attachments,
			ReadAt: func() *time.Time {
				if message.ReadAt.Valid {
					return &message.ReadAt.Time
				}
				return nil
			}(),
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		}

		log.Printf("Publishing AI response to Kafka: ID=%s, SessionID=%s, SenderType=%s, Message=%s",
			kafkaMessage.ID, kafkaMessage.SessionID, kafkaMessage.SenderType, kafkaMessage.Message)
		if err := h.kafkaService.PublishMessage(c.Context(), kafkaMessage); err != nil {
			// Log error, tapi tetap lanjutkan
			log.Printf("Failed to publish AI response to Kafka: %v", err)
		} else {
			log.Printf("Successfully published AI response to Kafka")
		}
	}

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "AI response sent successfully",
		Data:    response,
	})
}

// SendTypingIndicator godoc
// @Summary Send typing indicator from AI system
// @Description Endpoint for AI system to send typing indicator status
// @Tags ai
// @Accept json
// @Produce json
// @Param request body domain.AITypingRequest true "AI typing indicator status"
// @Success 200 {object} domain.ApiResponse
// @Failure 400 {object} domain.ApiResponse
// @Failure 401 {object} domain.ApiResponse
// @Security ApiKeyAuth
// @Router /api/ai/typing [post]
func (h *AIHandler) SendTypingIndicator(c *fiber.Ctx) error {
	var req domain.AITypingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if req.SessionID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ApiResponse{
			Success: false,
			Message: "Session ID is required",
			Error:   "validation failed",
		})
	}

	// Get user from context (AI user that was authenticated)
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.ApiResponse{
			Success: false,
			Message: "Authentication required",
			Error:   "user not found in context",
		})
	}

	// Use the provided user ID or default to authenticated user ID
	userID := user.ID
	if req.UserID != "" {
		userID = req.UserID
	}

	// Create typing indicator message for Kafka
	typingMessage := domain.TypingMessage{
		Type:      "typing_indicator",
		SessionID: req.SessionID,
		UserID:    userID,
		UserType:  "ai", // Set user type to AI
		IsTyping:  req.IsTyping,
		Timestamp: time.Now(),
	}

	// Publish typing indicator message to Kafka
	if err := h.kafkaService.PublishTypingIndicator(c.Context(), typingMessage); err != nil {
		log.Printf("Failed to publish AI typing indicator to Kafka: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ApiResponse{
			Success: false,
			Message: "Failed to send typing indicator",
			Error:   err.Error(),
		})
	}

	log.Printf("AI typing indicator sent: SessionID=%s, UserID=%s, IsTyping=%v",
		req.SessionID, userID, req.IsTyping)

	return c.JSON(domain.ApiResponse{
		Success: true,
		Message: "Typing indicator sent successfully",
	})
}
