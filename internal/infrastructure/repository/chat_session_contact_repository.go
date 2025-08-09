package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type chatSessionContactRepository struct {
	db *gorm.DB
}

func NewChatSessionContactRepository(db *gorm.DB) domain.ChatSessionContactRepository {
	return &chatSessionContactRepository{db: db}
}

func (r *chatSessionContactRepository) Create(ctx context.Context, contact *domain.ChatSessionContact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

func (r *chatSessionContactRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*domain.ChatSessionContact, error) {
	var contact domain.ChatSessionContact
	if err := r.db.WithContext(ctx).First(&contact, "session_id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}

func (r *chatSessionContactRepository) Update(ctx context.Context, contact *domain.ChatSessionContact) error {
	return r.db.WithContext(ctx).Save(contact).Error
}

func (r *chatSessionContactRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ChatSessionContact{}, "session_id = ?", sessionID).Error
}
