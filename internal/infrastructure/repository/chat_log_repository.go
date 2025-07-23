package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type chatLogRepository struct {
	db *gorm.DB
}

func NewChatLogRepository(db *gorm.DB) domain.ChatLogRepository {
	return &chatLogRepository{db: db}
}

func (r *chatLogRepository) Create(ctx context.Context, log *domain.ChatLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *chatLogRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.ChatLog, error) {
	var logs []*domain.ChatLog
	if err := r.db.WithContext(ctx).
		Preload("Session").
		Preload("User").
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *chatLogRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.ChatLog, error) {
	var logs []*domain.ChatLog
	if err := r.db.WithContext(ctx).
		Preload("Session").
		Preload("User").
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
