package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketEscalationRepository struct {
	db *gorm.DB
}

func NewTicketEscalationRepository(db *gorm.DB) domain.TicketEscalationRepository {
	return &ticketEscalationRepository{db: db}
}

func (r *ticketEscalationRepository) Create(ctx context.Context, escalation *domain.TicketEscalation) error {
	return r.db.WithContext(ctx).Create(escalation).Error
}

func (r *ticketEscalationRepository) GetByID(ctx context.Context, id string) (*domain.TicketEscalation, error) {
	var escalation domain.TicketEscalation
	err := r.db.WithContext(ctx).
		Preload("Ticket").
		Preload("FromDepartment").
		Preload("ToDepartment").
		Preload("FromAgent").
		Preload("ToAgent").
		Preload("EscalatedBy").
		First(&escalation, "id = ? AND deleted_at = 0", id).Error

	if err != nil {
		return nil, err
	}
	return &escalation, nil
}

func (r *ticketEscalationRepository) GetByTicketID(ctx context.Context, ticketID string) ([]domain.TicketEscalation, error) {
	var escalations []domain.TicketEscalation
	err := r.db.WithContext(ctx).
		Preload("FromDepartment").
		Preload("ToDepartment").
		Preload("FromAgent").
		Preload("ToAgent").
		Preload("EscalatedBy").
		Where("ticket_id = ? AND deleted_at = 0", ticketID).
		Order("escalated_at DESC").
		Find(&escalations).Error

	return escalations, err
}

func (r *ticketEscalationRepository) GetAll(ctx context.Context, filter *domain.TicketEscalationFilter) ([]domain.TicketEscalation, error) {
	var escalations []domain.TicketEscalation
	query := r.db.WithContext(ctx).Where("deleted_at = 0")

	// Apply filters
	if filter != nil {
		if filter.TicketID != "" {
			query = query.Where("ticket_id = ?", filter.TicketID)
		}
		if filter.FromLevel != nil {
			query = query.Where("from_level = ?", *filter.FromLevel)
		}
		if filter.ToLevel != nil {
			query = query.Where("to_level = ?", *filter.ToLevel)
		}
		if filter.FromDepartmentID != "" {
			query = query.Where("from_department_id = ?", filter.FromDepartmentID)
		}
		if filter.ToDepartmentID != "" {
			query = query.Where("to_department_id = ?", filter.ToDepartmentID)
		}
		if filter.EscalatedByID != "" {
			query = query.Where("escalated_by_id = ?", filter.EscalatedByID)
		}
		if filter.IsAutoEscalation != nil {
			query = query.Where("is_auto_escalation = ?", *filter.IsAutoEscalation)
		}
		if filter.TriggerType != "" {
			query = query.Where("trigger_type = ?", filter.TriggerType)
		}
		if filter.WasSuccessful != nil {
			query = query.Where("was_successful = ?", *filter.WasSuccessful)
		}
		if !filter.StartDate.IsZero() {
			query = query.Where("escalated_at >= ?", filter.StartDate)
		}
		if !filter.EndDate.IsZero() {
			query = query.Where("escalated_at <= ?", filter.EndDate)
		}

		// Pagination
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	err := query.
		Preload("Ticket").
		Preload("FromDepartment").
		Preload("ToDepartment").
		Preload("FromAgent").
		Preload("ToAgent").
		Preload("EscalatedBy").
		Order("escalated_at DESC").
		Find(&escalations).Error

	return escalations, err
}

func (r *ticketEscalationRepository) Update(ctx context.Context, escalation *domain.TicketEscalation) error {
	return r.db.WithContext(ctx).Save(escalation).Error
}

func (r *ticketEscalationRepository) UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&domain.TicketEscalation{}).
		Where("id = ? AND deleted_at = 0", id).
		Updates(updates).Error
}

func (r *ticketEscalationRepository) Delete(ctx context.Context, id string) error {
	now := time.Now().Unix()
	return r.db.WithContext(ctx).Model(&domain.TicketEscalation{}).
		Where("id = ? AND deleted_at = 0", id).
		Update("deleted_at", now).Error
}

func (r *ticketEscalationRepository) Count(ctx context.Context, filter *domain.TicketEscalationFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.TicketEscalation{}).Where("deleted_at = 0")

	// Apply same filters as GetAll
	if filter != nil {
		if filter.TicketID != "" {
			query = query.Where("ticket_id = ?", filter.TicketID)
		}
		if filter.FromLevel != nil {
			query = query.Where("from_level = ?", *filter.FromLevel)
		}
		if filter.ToLevel != nil {
			query = query.Where("to_level = ?", *filter.ToLevel)
		}
		if filter.FromDepartmentID != "" {
			query = query.Where("from_department_id = ?", filter.FromDepartmentID)
		}
		if filter.ToDepartmentID != "" {
			query = query.Where("to_department_id = ?", filter.ToDepartmentID)
		}
		if filter.EscalatedByID != "" {
			query = query.Where("escalated_by_id = ?", filter.EscalatedByID)
		}
		if filter.IsAutoEscalation != nil {
			query = query.Where("is_auto_escalation = ?", *filter.IsAutoEscalation)
		}
		if filter.TriggerType != "" {
			query = query.Where("trigger_type = ?", filter.TriggerType)
		}
		if filter.WasSuccessful != nil {
			query = query.Where("was_successful = ?", *filter.WasSuccessful)
		}
		if !filter.StartDate.IsZero() {
			query = query.Where("escalated_at >= ?", filter.StartDate)
		}
		if !filter.EndDate.IsZero() {
			query = query.Where("escalated_at <= ?", filter.EndDate)
		}
	}

	err := query.Count(&count).Error
	return count, err
}

// GetEscalationStatsByLevel returns escalation statistics grouped by level
func (r *ticketEscalationRepository) GetEscalationStatsByLevel(ctx context.Context, departmentID string, days int) ([]domain.EscalationLevelStats, error) {
	var stats []domain.EscalationLevelStats

	startDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT 
			to_level as level,
			COUNT(*) as total_escalations,
			COUNT(CASE WHEN was_successful = true THEN 1 END) as successful_escalations,
			COUNT(CASE WHEN is_auto_escalation = true THEN 1 END) as auto_escalations,
			AVG(EXTRACT(EPOCH FROM (resolved_at - escalated_at))/3600) as avg_resolution_hours
		FROM ticket_escalations 
		WHERE deleted_at = 0 
			AND escalated_at >= ?
			%s
		GROUP BY to_level 
		ORDER BY to_level
	`

	var args []interface{}
	args = append(args, startDate)

	if departmentID != "" {
		query = fmt.Sprintf(query, "AND to_department_id = ?")
		args = append(args, departmentID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&stats).Error
	return stats, err
}

// GetEscalationTrends returns escalation trends over time
func (r *ticketEscalationRepository) GetEscalationTrends(ctx context.Context, days int) ([]domain.EscalationTrend, error) {
	var trends []domain.EscalationTrend

	startDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT 
			DATE(escalated_at) as date,
			COUNT(*) as total_escalations,
			COUNT(CASE WHEN trigger_type = 'sla_breach' THEN 1 END) as sla_breach_escalations,
			COUNT(CASE WHEN trigger_type = 'manual' THEN 1 END) as manual_escalations,
			COUNT(CASE WHEN is_auto_escalation = true THEN 1 END) as auto_escalations
		FROM ticket_escalations 
		WHERE deleted_at = 0 
			AND escalated_at >= ?
		GROUP BY DATE(escalated_at) 
		ORDER BY DATE(escalated_at)
	`

	err := r.db.WithContext(ctx).Raw(query, startDate).Scan(&trends).Error
	return trends, err
}

// GetTopEscalationReasons returns most common escalation reasons
func (r *ticketEscalationRepository) GetTopEscalationReasons(ctx context.Context, limit int, days int) ([]domain.EscalationReason, error) {
	var reasons []domain.EscalationReason

	startDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT 
			trigger_type,
			COUNT(*) as count,
			(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER()) as percentage
		FROM ticket_escalations 
		WHERE deleted_at = 0 
			AND escalated_at >= ?
		GROUP BY trigger_type 
		ORDER BY count DESC
		LIMIT ?
	`

	err := r.db.WithContext(ctx).Raw(query, startDate, limit).Scan(&reasons).Error
	return reasons, err
}
