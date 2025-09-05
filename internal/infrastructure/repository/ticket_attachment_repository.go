package repository

import (
	"context"

	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketAttachmentRepository struct {
	db *gorm.DB
}

func NewTicketAttachmentRepository(db *gorm.DB) domain.TicketAttachmentRepository {
	return &ticketAttachmentRepository{db: db}
}

func (r *ticketAttachmentRepository) Create(ctx context.Context, attachment *domain.TicketAttachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *ticketAttachmentRepository) GetByID(ctx context.Context, id string) (*domain.TicketAttachment, error) {
	var attachment domain.TicketAttachment
	err := r.db.WithContext(ctx).
		First(&attachment, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *ticketAttachmentRepository) GetByTicketID(ctx context.Context, ticketID string) ([]*domain.TicketAttachment, error) {
	var attachments []*domain.TicketAttachment
	err := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID).
		Order("created_at ASC").
		Find(&attachments).Error

	return attachments, err
}

func (r *ticketAttachmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.TicketAttachment{}, "id = ?", id).Error
}
