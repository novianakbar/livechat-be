package repository

import (
	"context"
	"time"

	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketSLARepository struct {
	db *gorm.DB
}

func NewTicketSLARepository(db *gorm.DB) domain.TicketSLARepository {
	return &ticketSLARepository{db: db}
}

func (r *ticketSLARepository) Create(ctx context.Context, sla *domain.TicketSLA) error {
	return r.db.WithContext(ctx).Create(sla).Error
}

func (r *ticketSLARepository) GetByID(ctx context.Context, id string) (*domain.TicketSLA, error) {
	var sla domain.TicketSLA
	err := r.db.WithContext(ctx).
		First(&sla, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &sla, nil
}

func (r *ticketSLARepository) GetByTicketID(ctx context.Context, ticketID string) (*domain.TicketSLA, error) {
	var sla domain.TicketSLA
	err := r.db.WithContext(ctx).
		First(&sla, "ticket_id = ?", ticketID).Error

	if err != nil {
		return nil, err
	}
	return &sla, nil
}

func (r *ticketSLARepository) Update(ctx context.Context, sla *domain.TicketSLA) error {
	return r.db.WithContext(ctx).Save(sla).Error
}

func (r *ticketSLARepository) GetOverdueTickets(ctx context.Context) ([]*domain.TicketSLA, error) {
	var slas []*domain.TicketSLA
	now := time.Now()

	err := r.db.WithContext(ctx).
		Preload("Ticket").
		Where("(response_deadline < ? AND responded_at IS NULL) OR (resolution_deadline < ? AND resolved_at IS NULL)",
			now, now).
		Find(&slas).Error

	return slas, err
}

func (r *ticketSLARepository) GetBreachedTickets(ctx context.Context) ([]*domain.TicketSLA, error) {
	var slas []*domain.TicketSLA
	now := time.Now()

	err := r.db.WithContext(ctx).
		Preload("Ticket").
		Where("(response_deadline < ? AND responded_at IS NULL) OR (resolution_deadline < ? AND resolved_at IS NULL)",
			now, now).
		Find(&slas).Error

	return slas, err
}

func (r *ticketSLARepository) GetSLAMetrics(ctx context.Context, departmentID *string, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	query := r.db.WithContext(ctx).Model(&domain.TicketSLA{}).
		Joins("JOIN tickets ON ticket_sla.ticket_id = tickets.id").
		Where("tickets.created_at BETWEEN ? AND ?", dateFrom, dateTo)

	if departmentID != nil {
		query = query.Where("tickets.department_id = ?", *departmentID)
	}

	// Total tickets
	var totalTickets int64
	query.Count(&totalTickets)
	metrics["total_tickets"] = totalTickets

	// Response SLA compliance
	var responseCompliant int64
	query.Where("responded_at IS NOT NULL AND responded_at <= response_deadline").Count(&responseCompliant)

	// Resolution SLA compliance
	var resolutionCompliant int64
	query.Where("resolved_at IS NOT NULL AND resolved_at <= resolution_deadline").Count(&resolutionCompliant)

	if totalTickets > 0 {
		metrics["response_sla_compliance"] = float64(responseCompliant) / float64(totalTickets) * 100
		metrics["resolution_sla_compliance"] = float64(resolutionCompliant) / float64(totalTickets) * 100
	} else {
		metrics["response_sla_compliance"] = 0.0
		metrics["resolution_sla_compliance"] = 0.0
	}

	// Average response time (in hours)
	type avgResult struct {
		AvgHours float64
	}
	var responseAvg avgResult
	r.db.WithContext(ctx).Raw(`
		SELECT AVG(EXTRACT(EPOCH FROM (responded_at - tickets.created_at))/3600) as avg_hours
		FROM ticket_sla 
		JOIN tickets ON ticket_sla.ticket_id = tickets.id
		WHERE responded_at IS NOT NULL 
		AND tickets.created_at BETWEEN ? AND ?
		`+(func() string {
		if departmentID != nil {
			return " AND tickets.department_id = ?"
		}
		return ""
	}()),
		append([]interface{}{dateFrom, dateTo}, func() []interface{} {
			if departmentID != nil {
				return []interface{}{*departmentID}
			}
			return []interface{}{}
		}())...,
	).Scan(&responseAvg)

	metrics["avg_response_time_hours"] = responseAvg.AvgHours

	// Average resolution time (in hours)
	var resolutionAvg avgResult
	r.db.WithContext(ctx).Raw(`
		SELECT AVG(EXTRACT(EPOCH FROM (resolved_at - tickets.created_at))/3600) as avg_hours
		FROM ticket_sla 
		JOIN tickets ON ticket_sla.ticket_id = tickets.id
		WHERE resolved_at IS NOT NULL 
		AND tickets.created_at BETWEEN ? AND ?
		`+(func() string {
		if departmentID != nil {
			return " AND tickets.department_id = ?"
		}
		return ""
	}()),
		append([]interface{}{dateFrom, dateTo}, func() []interface{} {
			if departmentID != nil {
				return []interface{}{*departmentID}
			}
			return []interface{}{}
		}())...,
	).Scan(&resolutionAvg)

	metrics["avg_resolution_time_hours"] = resolutionAvg.AvgHours

	return metrics, nil
}
