package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type chatSessionRepository struct {
	db *gorm.DB
}

func NewChatSessionRepository(db *gorm.DB) domain.ChatSessionRepository {
	return &chatSessionRepository{db: db}
}

func (r *chatSessionRepository) Create(ctx context.Context, session *domain.ChatSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *chatSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ChatSession, error) {
	var session domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		First(&session, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *chatSessionRepository) GetByChatUserID(ctx context.Context, chatUserID uuid.UUID) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("chat_user_id = ?", chatUserID).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) GetByAgentID(ctx context.Context, agentID uuid.UUID) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) GetActiveSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("status = ?", "active").
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) GetWaitingSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("status = ?", "waiting").
		Order("created_at ASC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) Update(ctx context.Context, session *domain.ChatSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *chatSessionRepository) Close(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.ChatSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":   "closed",
			"ended_at": &now,
		}).Error
}

func (r *chatSessionRepository) GetSessionsByStatus(ctx context.Context, status string) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) GetSessionsByDateRange(ctx context.Context, start, end time.Time) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// Analytics methods
func (r *chatSessionRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.ChatSession{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *chatSessionRepository) CountCompletedSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.ChatSession{}).
		Where("status = ? AND ended_at >= ?", "closed", since).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *chatSessionRepository) GetAverageResponseTime(ctx context.Context) (float64, error) {
	// This would calculate average time from session start to first agent response
	// For now, return a mock value (in seconds)
	return 180.0, nil // 3 minutes average
}

func (r *chatSessionRepository) GetOSSCategoriesStats(ctx context.Context) ([]domain.CategoryStats, error) {
	// This would analyze topics to categorize OSS requests
	// For now, return mock data
	return []domain.CategoryStats{
		{Category: "NIB (Nomor Induk Berusaha)", Count: 345, Percentage: 35},
		{Category: "Izin Usaha Perdagangan", Count: 289, Percentage: 29},
		{Category: "Izin Usaha Industri", Count: 198, Percentage: 20},
		{Category: "Izin Usaha Jasa", Count: 156, Percentage: 16},
	}, nil
}

func (r *chatSessionRepository) GetWithPagination(ctx context.Context, offset, limit int, status string, agentID, departmentID *uuid.UUID) ([]*domain.ChatSession, error) {
	query := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if agentID != nil {
		query = query.Where("agent_id = ?", *agentID)
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	var sessions []*domain.ChatSession
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *chatSessionRepository) Count(ctx context.Context, status string, agentID, departmentID *uuid.UUID) (int, error) {
	query := r.db.WithContext(ctx).Model(&domain.ChatSession{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if agentID != nil {
		query = query.Where("agent_id = ?", *agentID)
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *chatSessionRepository) GetSessionsWithMessages(ctx context.Context, chatUserID uuid.UUID, limit, offset int) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("chat_user_id = ?", chatUserID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) GetSessionHistory(ctx context.Context, chatUserID uuid.UUID, limit, offset int) ([]*domain.ChatSession, error) {
	var sessions []*domain.ChatSession
	if err := r.db.WithContext(ctx).
		Preload("ChatUser").
		Preload("Agent").
		Preload("Department").
		Preload("Contact").
		Where("chat_user_id = ?", chatUserID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
