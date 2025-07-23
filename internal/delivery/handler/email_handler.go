package handler

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
)

type EmailHandler struct {
	EmailUseCase domain.EmailService
}

func NewEmailHandler(emailUseCase domain.EmailService) *EmailHandler {
	return &EmailHandler{
		EmailUseCase: emailUseCase,
	}
}

// SendEmail handles sending a custom email
func (h *EmailHandler) SendEmail(c *fiber.Ctx) error {
	var req domain.SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}
	resp, err := h.EmailUseCase.SendEmail(context.Background(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}
	return c.JSON(resp)
}

// SendWelcomeEmail handles sending a welcome email
func (h *EmailHandler) SendWelcomeEmail(c *fiber.Ctx) error {
	type welcomeReq struct {
		To   string `json:"to"`
		Name string `json:"name"`
	}
	var req welcomeReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}
	resp, err := h.EmailUseCase.SendWelcomeEmail(context.Background(), req.To, req.Name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}
	return c.JSON(resp)
}

// SendPasswordResetEmail handles sending a password reset email
func (h *EmailHandler) SendPasswordResetEmail(c *fiber.Ctx) error {
	type resetReq struct {
		To         string `json:"to"`
		ResetToken string `json:"reset_token"`
	}
	var req resetReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}
	resp, err := h.EmailUseCase.SendPasswordResetEmail(context.Background(), req.To, req.ResetToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}
	return c.JSON(resp)
}

// SendChatTranscriptEmail handles sending a chat transcript email
func (h *EmailHandler) SendChatTranscriptEmail(c *fiber.Ctx) error {
	type transcriptReq struct {
		To         string `json:"to"`
		Transcript string `json:"transcript"`
		SessionID  string `json:"session_id"`
	}
	var req transcriptReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}
	sessionUUID, err := uuid.Parse(req.SessionID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid session_id", "error": err.Error()})
	}
	resp, err := h.EmailUseCase.SendChatTranscriptEmail(context.Background(), req.To, req.Transcript, sessionUUID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}
	return c.JSON(resp)
}

// SendCustomEmail handles sending a custom template email
func (h *EmailHandler) SendCustomEmail(c *fiber.Ctx) error {
	type customReq struct {
		Template  domain.EmailTemplate `json:"template"`
		To        []string             `json:"to"`
		Variables map[string]string    `json:"variables"`
	}
	var req customReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request", "error": err.Error()})
	}
	resp, err := h.EmailUseCase.SendTemplatedEmail(context.Background(), &req.Template, req.To, req.Variables)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}
	return c.JSON(resp)
}
