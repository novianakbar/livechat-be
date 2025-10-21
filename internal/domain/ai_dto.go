package domain

import "github.com/google/uuid"

type AITypingRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	IsTyping  bool      `json:"is_typing"`
	UserID    string    `json:"user_id,omitempty"` // Opsional, default ke ID AI yang dikonfigurasi
}
