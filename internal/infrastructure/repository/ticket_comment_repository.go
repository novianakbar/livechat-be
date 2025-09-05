package repository

import (
	"context"

	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketCommentRepository struct {
	db *gorm.DB
}

func NewTicketCommentRepository(db *gorm.DB) domain.TicketCommentRepository {
	return &ticketCommentRepository{db: db}
}

func (r *ticketCommentRepository) Create(ctx context.Context, comment *domain.TicketComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *ticketCommentRepository) GetByID(ctx context.Context, id string) (*domain.TicketComment, error) {
	var comment domain.TicketComment
	err := r.db.WithContext(ctx).
		Preload("CreatedBy").
		First(&comment, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *ticketCommentRepository) GetByTicketID(ctx context.Context, ticketID string, includeInternal bool) ([]*domain.TicketComment, error) {
	var comments []*domain.TicketComment
	query := r.db.WithContext(ctx).
		Preload("CreatedBy").
		Where("ticket_id = ?", ticketID)

	if !includeInternal {
		query = query.Where("is_public = ?", true)
	}

	err := query.Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *ticketCommentRepository) GetPublicByTicketID(ctx context.Context, ticketID string) ([]*domain.TicketComment, error) {
	var comments []*domain.TicketComment
	err := r.db.WithContext(ctx).
		Preload("CreatedBy").
		Where("ticket_id = ? AND is_public = ?", ticketID, true).
		Order("created_at ASC").
		Find(&comments).Error

	return comments, err
}

func (r *ticketCommentRepository) Update(ctx context.Context, comment *domain.TicketComment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *ticketCommentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.TicketComment{}, "id = ?", id).Error
}
