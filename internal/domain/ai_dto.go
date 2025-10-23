package domain

import "github.com/google/uuid"

type AITypingRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	IsTyping  bool      `json:"is_typing"`
	UserID    string    `json:"user_id,omitempty"` // Opsional, default ke ID AI yang dikonfigurasi
}

type AIResponseRequest struct {
	SessionID              uuid.UUID `json:"session_id" validate:"required"`
	Message                string    `json:"message" validate:"required"`
	MessageType            string    `json:"message_type" validate:"oneof=text image file system"`
	Attachments            []string  `json:"attachments"`
	OfferHumanEscalation   bool      `json:"offer_human_escalation"`   // Flag untuk menampilkan penawaran eskalasi
	EscalationOfferMessage string    `json:"escalation_offer_message"` // Pesan penawaran eskalasi (opsional)
}

type CustomerEscalationRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	Accept    bool      `json:"accept" validate:"required"` // true = terima eskalasi, false = tolak
}
