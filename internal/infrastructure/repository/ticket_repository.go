package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) domain.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	// Generate ticket code if not provided
	if ticket.TicketCode == "" {
		ticket.TicketCode = r.generateTicketCode()
	}

	// Generate access token if not provided
	if ticket.AccessToken == "" {
		ticket.AccessToken = r.generateAccessToken()
	}

	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) GetByID(ctx context.Context, id string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Attachments").
		Preload("History").
		Preload("SLA").
		First(&ticket, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByTicketCode(ctx context.Context, code string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Attachments").
		Preload("History").
		Preload("SLA").
		First(&ticket, "ticket_code = ?", code).Error

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetByAccessToken(ctx context.Context, token string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Preload("Comments").
		Preload("Attachments").
		Preload("History").
		Preload("SLA").
		First(&ticket, "access_token = ?", token).Error

	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket) error {
	// Use Select to only update specific fields and avoid foreign key issues with preloaded relations
	// Exclude fields that should not be updated: created_via, ticket_code, access_token, created_at, etc.
	return r.db.WithContext(ctx).Model(ticket).Select(
		"subject",
		"description",
		"customer_name",
		"customer_email",
		"customer_phone",
		"category_id",
		"priority",
		"status",
		"assigned_to_id",
		"department_id",
		"first_response_at",
		"resolved_at",
		"closed_at",
		"updated_at",
	).Updates(ticket).Error
}

// UpdateFields - Safe method for updating specific fields without constraint violations
func (r *ticketRepository) UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error {
	// Always add updated_at
	updates["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Ticket{}, "id = ?", id).Error
}

func (r *ticketRepository) GetList(ctx context.Context, filter *domain.TicketFilter) ([]*domain.Ticket, int64, error) {
	var tickets []*domain.Ticket
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Ticket{})

	// Apply filters
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedTo)
	}
	if filter.CreatedBy != nil {
		query = query.Where("created_by = ?", *filter.CreatedBy)
	}
	if filter.CustomerInfo != nil {
		searchTerm := "%" + *filter.CustomerInfo + "%"
		query = query.Where("customer_name ILIKE ? OR customer_email ILIKE ? OR customer_phone ILIKE ?",
			searchTerm, searchTerm, searchTerm)
	}
	if filter.Subject != nil {
		query = query.Where("subject ILIKE ?", "%"+*filter.Subject+"%")
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := "created_at desc" // default
	if filter.SortBy != "" {
		direction := "desc"
		if filter.SortOrder == "asc" {
			direction = "asc"
		}
		orderBy = fmt.Sprintf("%s %s", filter.SortBy, direction)
	}
	query = query.Order(orderBy)

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Load with preloads
	err := query.
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Find(&tickets).Error

	return tickets, total, err
}

func (r *ticketRepository) GetByAssignedAgent(ctx context.Context, agentID string, status []string) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	query := r.db.WithContext(ctx).Where("assigned_to = ?", agentID)

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Order("created_at desc").
		Find(&tickets).Error

	return tickets, err
}

func (r *ticketRepository) GetByDepartment(ctx context.Context, deptID string, status []string) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	query := r.db.WithContext(ctx).Where("department_id = ?", deptID)

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Order("created_at desc").
		Find(&tickets).Error

	return tickets, err
}

func (r *ticketRepository) GetOverdueTickets(ctx context.Context) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket

	// Get tickets that are overdue based on SLA
	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN ticket_sla ON tickets.id = ticket_sla.ticket_id").
		Where("tickets.status NOT IN ?", []string{"closed", "resolved"}).
		Where("(ticket_sla.response_deadline < ? AND ticket_sla.responded_at IS NULL) OR (ticket_sla.resolution_deadline < ? AND ticket_sla.resolved_at IS NULL)",
			time.Now(), time.Now()).
		Preload("Category").
		Preload("Department").
		Preload("AssignedTo").
		Preload("CreatedBy").
		Preload("SLA").
		Find(&tickets).Error

	return tickets, err
}

// Helper functions
func (r *ticketRepository) generateTicketCode() string {
	// Format: TKT-YYYYMMDD-XXXXX
	now := time.Now()
	dateStr := now.Format("20060102")

	// Get count of tickets created today
	var count int64
	r.db.Model(&domain.Ticket{}).
		Where("DATE(created_at) = ?", now.Format("2006-01-02")).
		Count(&count)

	return fmt.Sprintf("TKT-%s-%05d", dateStr, count+1)
}

func (r *ticketRepository) generateAccessToken() string {
	// Generate a unique token for customer access
	token := strings.ReplaceAll(uuid.New().String(), "-", "")
	return strings.ToUpper(token[:16]) // 16 character token
}

// ============================================
// STATISTICS METHODS IMPLEMENTATION
// ============================================

// CountByStatus counts tickets by status
func (r *ticketRepository) CountByStatus(ctx context.Context, statuses []string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.Ticket{})

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	err := query.Count(&count).Error
	return count, err
}

// CountEscalatedTickets counts tickets that have been escalated
func (r *ticketRepository) CountEscalatedTickets(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("escalation_count > 0 OR current_level > 0").
		Count(&count).Error
	return count, err
}

// CountOverdueTickets counts tickets that are overdue (SLA breach)
func (r *ticketRepository) CountOverdueTickets(ctx context.Context) (int64, error) {
	var count int64
	now := time.Now()

	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Joins("LEFT JOIN ticket_slas ON tickets.id = ticket_slas.ticket_id").
		Where("(ticket_slas.first_response_due < ? AND tickets.first_response_at IS NULL) OR "+
			"(ticket_slas.resolution_due < ? AND ticket_slas.resolved_at IS NULL)", now, now).
		Count(&count).Error
	return count, err
}

// CountCreatedSince counts tickets created since a specific time
func (r *ticketRepository) CountCreatedSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("created_at >= ?", since).
		Count(&count).Error
	return count, err
}

// CountResolvedSince counts tickets resolved since a specific time
func (r *ticketRepository) CountResolvedSince(ctx context.Context, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Where("resolved_at >= ?", since).
		Count(&count).Error
	return count, err
}

// GetAverageResolutionTime calculates average resolution time in hours
func (r *ticketRepository) GetAverageResolutionTime(ctx context.Context) (float64, error) {
	var result struct {
		AvgTime float64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT AVG(EXTRACT(EPOCH FROM (resolved_at - created_at))/3600) as avg_time 
		FROM tickets 
		WHERE resolved_at IS NOT NULL
	`).Scan(&result).Error

	return result.AvgTime, err
}

// GetAverageFirstResponseTime calculates average first response time in hours
func (r *ticketRepository) GetAverageFirstResponseTime(ctx context.Context) (float64, error) {
	var result struct {
		AvgTime float64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT AVG(EXTRACT(EPOCH FROM (first_response_at - created_at))/3600) as avg_time 
		FROM tickets 
		WHERE first_response_at IS NOT NULL
	`).Scan(&result).Error

	return result.AvgTime, err
}

// GetSLAComplianceRate calculates SLA compliance rate as percentage
func (r *ticketRepository) GetSLAComplianceRate(ctx context.Context) (float64, error) {
	var result struct {
		ComplianceRate float64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COALESCE(
				(COUNT(CASE WHEN NOT first_response_breached AND NOT resolution_breached THEN 1 END) * 100.0 / 
				 NULLIF(COUNT(*), 0)), 
				0
			) as compliance_rate
		FROM ticket_slas
	`).Scan(&result).Error

	return result.ComplianceRate, err
}

// GetStatusBreakdown returns count of tickets by status
func (r *ticketRepository) GetStatusBreakdown(ctx context.Context) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	breakdown := make(map[string]int)
	for _, result := range results {
		breakdown[result.Status] = result.Count
	}

	return breakdown, nil
}

// GetPriorityBreakdown returns count of tickets by priority
func (r *ticketRepository) GetPriorityBreakdown(ctx context.Context) (map[string]int, error) {
	var results []struct {
		Priority string
		Count    int
	}

	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Select("priority, COUNT(*) as count").
		Group("priority").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	breakdown := make(map[string]int)
	for _, result := range results {
		breakdown[result.Priority] = result.Count
	}

	return breakdown, nil
}

// GetCategoryBreakdown returns count of tickets by category
func (r *ticketRepository) GetCategoryBreakdown(ctx context.Context) (map[string]int, error) {
	var results []struct {
		CategoryName string
		Count        int
	}

	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Select("ticket_categories.name as category_name, COUNT(*) as count").
		Joins("LEFT JOIN ticket_categories ON tickets.category_id = ticket_categories.id").
		Group("ticket_categories.name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	breakdown := make(map[string]int)
	for _, result := range results {
		breakdown[result.CategoryName] = result.Count
	}

	return breakdown, nil
}

// GetLevelBreakdown returns count of tickets by escalation level
func (r *ticketRepository) GetLevelBreakdown(ctx context.Context) (map[string]int, error) {
	var results []struct {
		Level string
		Count int
	}

	err := r.db.WithContext(ctx).Model(&domain.Ticket{}).
		Select("CONCAT('L', current_level) as level, COUNT(*) as count").
		Group("current_level").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	breakdown := make(map[string]int)
	for _, result := range results {
		breakdown[result.Level] = result.Count
	}

	return breakdown, nil
}
