package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) domain.ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) Create(ctx context.Context, message *domain.ChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *chatMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error) {
	var message domain.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Session").
		Preload("Sender").
		First(&message, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *chatMessageRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChatMessage, error) {
	var messages []*domain.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Session").
		Preload("Sender").
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatMessageRepository) Update(ctx context.Context, message *domain.ChatMessage) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *chatMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ChatMessage{}, "id = ?", id).Error
}

func (r *chatMessageRepository) MarkAsRead(ctx context.Context, messageID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.ChatMessage{}).
		Where("id = ?", messageID).
		Update("read_at", "NOW()").Error
}

func (r *chatMessageRepository) GetUnreadMessages(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChatMessage, error) {
	var messages []*domain.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Session").
		Preload("Sender").
		Where("session_id = ? AND read_at IS NULL", sessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatMessageRepository) GetMessagesByDateRange(ctx context.Context, start, end time.Time) ([]*domain.ChatMessage, error) {
	var messages []*domain.ChatMessage
	if err := r.db.WithContext(ctx).
		Preload("Session").
		Preload("Sender").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// Analytics methods
func (r *chatMessageRepository) GetTopQuestions(ctx context.Context, limit int) ([]domain.QuestionStats, error) {
	// This would analyze message content to find common questions
	// For now, return mock data based on OSS common questions
	return []domain.QuestionStats{
		{Question: "Bagaimana cara mengurus NIB?", Count: 234},
		{Question: "Berapa lama proses penerbitan izin?", Count: 187},
		{Question: "Dokumen apa saja yang diperlukan?", Count: 156},
		{Question: "Apakah ada biaya pengurusan?", Count: 142},
		{Question: "Status permohonan saya bagaimana?", Count: 98},
	}, nil
}
