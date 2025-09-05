package repository

import (
	"context"

	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketHistoryRepository struct {
	db *gorm.DB
}

func NewTicketHistoryRepository(db *gorm.DB) domain.TicketHistoryRepository {
	return &ticketHistoryRepository{db: db}
}

func (r *ticketHistoryRepository) Create(ctx context.Context, history *domain.TicketHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *ticketHistoryRepository) GetByID(ctx context.Context, id string) (*domain.TicketHistory, error) {
	var history domain.TicketHistory
	err := r.db.WithContext(ctx).
		Preload("ChangedByUser").
		First(&history, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *ticketHistoryRepository) GetByTicketID(ctx context.Context, ticketID string) ([]*domain.TicketHistory, error) {
	var histories []*domain.TicketHistory
	err := r.db.WithContext(ctx).
		Preload("ChangedByUser").
		Where("ticket_id = ?", ticketID).
		Order("created_at ASC").
		Find(&histories).Error

	return histories, err
}
