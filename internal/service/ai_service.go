package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/pkg/config"
)

type AIWebhookRequest struct {
	ChatID    string `json:"chat_id"`
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
}

type aiService struct {
	config *config.Config
}

func NewAIService(config *config.Config) domain.AIService {
	return &aiService{
		config: config,
	}
}

func (s *aiService) SendToAI(ctx context.Context, message *domain.ChatMessage) error {
	if !s.config.AI.Enabled || s.config.AI.WebhookURL == "" {
		return errors.New("AI integration is not enabled or webhook URL is not configured")
	}

	// Prepare request body
	userID := ""
	if message.SenderID.Valid {
		userID = message.SenderID.String
	}

	req := AIWebhookRequest{
		ChatID:    message.ID,
		SessionID: message.SessionID,
		UserID:    userID,
		Message:   message.Message,
	}

	// Convert to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.AI.WebhookURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Send request
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return errors.New("AI webhook returned non-success status code: " + resp.Status)
	}

	return nil
}
