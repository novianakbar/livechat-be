package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/pkg/utils"
)

type ticketService struct {
	ticketRepo     domain.TicketRepository
	categoryRepo   domain.TicketCategoryRepository
	commentRepo    domain.TicketCommentRepository
	historyRepo    domain.TicketHistoryRepository
	slaRepo        domain.TicketSLARepository
	escalationRepo domain.TicketEscalationRepository
	userRepo       domain.UserRepository
	departmentRepo domain.DepartmentRepository
}

func NewTicketService(
	ticketRepo domain.TicketRepository,
	categoryRepo domain.TicketCategoryRepository,
	commentRepo domain.TicketCommentRepository,
	historyRepo domain.TicketHistoryRepository,
	slaRepo domain.TicketSLARepository,
	escalationRepo domain.TicketEscalationRepository,
	userRepo domain.UserRepository,
	departmentRepo domain.DepartmentRepository,
) domain.TicketService {
	return &ticketService{
		ticketRepo:     ticketRepo,
		categoryRepo:   categoryRepo,
		commentRepo:    commentRepo,
		historyRepo:    historyRepo,
		slaRepo:        slaRepo,
		escalationRepo: escalationRepo,
		userRepo:       userRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *ticketService) CreateTicket(ctx context.Context, req *domain.CreateTicketRequest) (*domain.Ticket, error) {
	// Validate category exists
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Create ticket
	ticket := &domain.Ticket{
		Subject:       req.Subject,
		Description:   req.Description,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		Priority:      req.Priority,
		Status:        "open",
		CreatedVia:    req.CreatedVia,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set category ID
	ticket.CategoryID = sql.NullString{String: req.CategoryID, Valid: true}

	// Set department ID if category has one
	if category.DefaultDepartmentID.Valid {
		ticket.DepartmentID = category.DefaultDepartmentID
	}

	// Set created by if provided
	if req.CreatedBy != nil {
		ticket.CreatedByID = sql.NullString{String: *req.CreatedBy, Valid: true}
	}

	err = s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Create SLA record
	err = s.createSLARecord(ctx, ticket, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create SLA record: %w", err)
	}

	return s.ticketRepo.GetByID(ctx, ticket.ID)
}

func (s *ticketService) GetTicket(ctx context.Context, id string) (*domain.Ticket, error) {
	return s.ticketRepo.GetByID(ctx, id)
}

func (s *ticketService) GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error) {
	return s.ticketRepo.GetByTicketCode(ctx, code)
}

func (s *ticketService) GetTicketByAccessToken(ctx context.Context, token string) (*domain.Ticket, error) {
	return s.ticketRepo.GetByAccessToken(ctx, token)
}

func (s *ticketService) UpdateTicket(ctx context.Context, req *domain.UpdateTicketRequest) (*domain.Ticket, error) {
	// First check if ticket exists
	_, err := s.ticketRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Convert request struct to update map, automatically excluding non-updatable fields
	updates := utils.ToUpdateMap(req)

	// Only update if there are fields to update
	if len(updates) > 1 { // > 1 because updated_at is always added
		err = s.ticketRepo.UpdateFields(ctx, req.ID, updates)
		if err != nil {
			return nil, fmt.Errorf("failed to update ticket: %w", err)
		}
	}

	return s.ticketRepo.GetByID(ctx, req.ID)
}

func (s *ticketService) AssignTicket(ctx context.Context, req *domain.AssignTicketRequest) error {
	// Check if ticket exists
	_, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Validate agent exists
	_, err = s.userRepo.GetByID(ctx, req.AgentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Update only the assigned_to_id field
	updates := map[string]interface{}{
		"assigned_to_id": req.AgentID,
	}

	return s.ticketRepo.UpdateFields(ctx, req.TicketID, updates)
}

func (s *ticketService) EscalateTicket(ctx context.Context, req *domain.EscalateTicketRequest) error {
	// 1. Get current ticket with full details
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// 2. Determine next escalation level
	nextLevel := ticket.CurrentLevel + 1
	if req.ToLevel != nil {
		nextLevel = *req.ToLevel
	}

	// 3. Validate escalation is possible
	if nextLevel > ticket.MaxLevel {
		return fmt.Errorf("cannot escalate beyond maximum level %d", ticket.MaxLevel)
	}

	if !ticket.CanEscalate {
		return fmt.Errorf("ticket cannot be escalated")
	}

	// 4. Find appropriate department for next level
	targetDept, err := s.GetEscalationPath(ctx, ticket.DepartmentID.String, nextLevel)
	if err != nil {
		return fmt.Errorf("no suitable escalation department: %w", err)
	}

	// 5. Auto-assign to appropriate agent if not specified
	var targetAgent *domain.User
	if req.AssignToID != nil {
		targetAgent, err = s.userRepo.GetByID(ctx, *req.AssignToID)
		if err != nil {
			return fmt.Errorf("specified agent not found: %w", err)
		}
		// Validate agent belongs to target department
		if targetAgent.DepartmentID.String != targetDept.ID {
			return fmt.Errorf("agent does not belong to target department")
		}
	} else {
		targetAgent, err = s.GetBestAvailableAgent(ctx, targetDept.ID, nextLevel)
		if err != nil {
			return fmt.Errorf("failed to assign agent: %w", err)
		}
	}

	// 6. Record escalation history
	escalationRecord := &domain.TicketEscalation{
		TicketID:         req.TicketID,
		FromLevel:        ticket.CurrentLevel,
		ToLevel:          nextLevel,
		FromDepartmentID: ticket.DepartmentID,
		ToDepartmentID:   sql.NullString{String: targetDept.ID, Valid: true},
		FromAgentID:      ticket.AssignedToID,
		ToAgentID:        sql.NullString{String: targetAgent.ID, Valid: true},
		Reason:           req.Reason,
		EscalatedByID:    req.EscalatedBy,
		IsAutoEscalation: false, // This is manual escalation
		TriggerType:      "manual",
		WasSuccessful:    true,
	}

	if err := s.escalationRepo.Create(ctx, escalationRecord); err != nil {
		return fmt.Errorf("failed to record escalation: %w", err)
	}

	// 7. Update ticket with new level and assignment
	updates := map[string]interface{}{
		"current_level":           nextLevel,
		"status":                  "escalated",
		"previous_assigned_to_id": ticket.AssignedToID,
		"previous_department_id":  ticket.DepartmentID,
		"assigned_to_id":          targetAgent.ID,
		"department_id":           targetDept.ID,
		"escalation_count":        ticket.EscalationCount + 1,
	}

	if err := s.ticketRepo.UpdateFields(ctx, req.TicketID, updates); err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	// 8. Create history record
	historyDescription := fmt.Sprintf("Escalated from L%d (%s) to L%d (%s): %s",
		ticket.CurrentLevel,
		getDepName(ticket.Department),
		nextLevel,
		targetDept.Name,
		req.Reason)

	if err := s.CreateHistoryRecord(ctx, req.TicketID, "escalated", historyDescription, req.EscalatedBy); err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	// 9. Send escalation notification
	go s.sendEscalationNotification(ctx, ticket, ticket.CurrentLevel, nextLevel)

	// 10. Check and update SLA if needed
	if ticket.SLA != nil {
		// SLA may need adjustment based on new level
		if err := s.adjustSLAForEscalation(ctx, ticket, nextLevel); err != nil {
			// Log error but don't fail the escalation
			fmt.Printf("Warning: Failed to adjust SLA for escalation: %v\n", err)
		}
	}

	return nil
}

func (s *ticketService) AddComment(ctx context.Context, req *domain.AddCommentRequest) (*domain.TicketComment, error) {
	// Get ticket to validate
	_, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Check for first response from agent to customer (only for public comments)
	if req.IsPublic {
		comments, err := s.commentRepo.GetByTicketID(ctx, req.TicketID, false) // Get only public comments
		if err != nil {
			return nil, fmt.Errorf("failed to get existing comments: %w", err)
		}

		// Check if this is the first public response from an agent
		isFirstResponse := true
		for _, comment := range comments {
			// If there's already a public comment from an agent (not customer), not first response
			if !comment.IsFromCustomer {
				isFirstResponse = false
				break
			}
		}

		// Update first response time if this is the first agent response
		if isFirstResponse {
			err := s.UpdateFirstResponse(ctx, req.TicketID)
			if err != nil {
				return nil, fmt.Errorf("failed to update first response time: %w", err)
			}
		}
	}

	// Create comment
	comment := &domain.TicketComment{
		TicketID:       req.TicketID,
		UserID:         sql.NullString{String: req.CreatedBy, Valid: true},
		Content:        req.Content,
		IsInternal:     !req.IsPublic,
		IsFromCustomer: false, // Assuming agent comment; adjust if needed
		AuthorName:     "",    // Will be populated from user data
		AuthorEmail:    "",    // Will be populated from user data
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Update ticket's updated_at timestamp
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	err = s.ticketRepo.UpdateFields(ctx, req.TicketID, updates)
	if err != nil {
		// Log error but don't fail the comment creation
		fmt.Printf("Warning: failed to update ticket timestamp: %v\n", err)
	}

	return comment, nil
}

func (s *ticketService) GetTicketList(ctx context.Context, filter *domain.TicketFilter) ([]*domain.Ticket, int64, error) {
	return s.ticketRepo.GetList(ctx, filter)
}

func (s *ticketService) GetAgentTickets(ctx context.Context, agentID string, status []string) ([]*domain.Ticket, error) {
	return s.ticketRepo.GetByAssignedAgent(ctx, agentID, status)
}

func (s *ticketService) GetDepartmentTickets(ctx context.Context, deptID string, status []string) ([]*domain.Ticket, error) {
	return s.ticketRepo.GetByDepartment(ctx, deptID, status)
}

// Helper function to create SLA record
func (s *ticketService) createSLARecord(ctx context.Context, ticket *domain.Ticket, category *domain.TicketCategory) error {
	now := time.Now()

	sla := &domain.TicketSLA{
		TicketID:         ticket.ID,
		FirstResponseDue: now.Add(time.Duration(category.SLAFirstResponse) * time.Minute),
		ResolutionDue:    now.Add(time.Duration(category.SLAResolution) * time.Minute),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return s.slaRepo.Create(ctx, sla)
}

// Helper methods for safe updates
func (s *ticketService) UpdateTicketStatus(ctx context.Context, ticketID, status string) error {
	// Get current ticket state
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	oldStatus := ticket.Status

	// Validate status transition
	if !s.isValidStatusTransition(ticket.Status, status) {
		return fmt.Errorf("invalid status transition from %s to %s", ticket.Status, status)
	}

	// Prepare updates
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// Add timestamp for specific status changes
	now := time.Now()
	switch status {
	case "resolved":
		updates["resolved_at"] = now
	case "closed":
		updates["closed_at"] = now
	case "in_progress":
		// Ensure ticket is assigned when moved to in_progress
		if !ticket.AssignedToID.Valid {
			return fmt.Errorf("ticket must be assigned before setting to in_progress")
		}
	}

	err = s.ticketRepo.UpdateFields(ctx, ticketID, updates)
	if err != nil {
		return err
	}

	// Send status change notification
	go s.sendStatusChangeNotification(ctx, ticket, oldStatus, status)

	return nil
}

// isValidStatusTransition validates business rules for status changes
func (s *ticketService) isValidStatusTransition(currentStatus, newStatus string) bool {
	// Define valid transitions
	validTransitions := map[string][]string{
		"open":        {"in_progress", "escalated", "closed"},
		"in_progress": {"resolved", "escalated", "open", "closed"},
		"resolved":    {"closed", "open"}, // Can reopen resolved tickets
		"escalated":   {"in_progress", "resolved", "closed"},
		"closed":      {"open"}, // Can reopen closed tickets with proper authorization
	}

	allowedStates, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedStates {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// sendCustomerNotification sends notification to customer (placeholder implementation)
func (s *ticketService) sendCustomerNotification(ctx context.Context, ticket *domain.Ticket, comment *domain.TicketComment) {
	// TODO: Implement email/SMS notification system
	// For now, just log the notification
	fmt.Printf("NOTIFICATION: Ticket %s received new comment from agent\n", ticket.TicketCode)
	fmt.Printf("Customer: %s (%s)\n", ticket.CustomerName, ticket.CustomerEmail)
	fmt.Printf("Comment: %s\n", comment.Content[:min(50, len(comment.Content))])
}

// sendStatusChangeNotification sends notification for status changes
func (s *ticketService) sendStatusChangeNotification(ctx context.Context, ticket *domain.Ticket, oldStatus, newStatus string) {
	// TODO: Implement email/SMS notification system
	// For now, just log the notification
	fmt.Printf("NOTIFICATION: Ticket %s status changed from %s to %s\n", ticket.TicketCode, oldStatus, newStatus)
	fmt.Printf("Customer: %s (%s)\n", ticket.CustomerName, ticket.CustomerEmail)
}

// sendEscalationNotification sends notification for escalations
func (s *ticketService) sendEscalationNotification(ctx context.Context, ticket *domain.Ticket, fromLevel, toLevel int) {
	// TODO: Implement email/SMS notification system
	// For now, just log the notification
	fmt.Printf("NOTIFICATION: Ticket %s escalated from L%d to L%d\n", ticket.TicketCode, fromLevel, toLevel)
	if ticket.AssignedTo != nil {
		fmt.Printf("Assigned to: %s (%s)\n", ticket.AssignedTo.Name, ticket.AssignedTo.Email)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *ticketService) UpdateFirstResponse(ctx context.Context, ticketID string) error {
	updates := map[string]interface{}{
		"first_response_at": time.Now(),
	}

	return s.ticketRepo.UpdateFields(ctx, ticketID, updates)
}

func (s *ticketService) UnassignTicket(ctx context.Context, ticketID string) error {
	// Check if ticket exists
	_, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Remove assignment
	updates := map[string]interface{}{
		"assigned_to_id": nil,
	}

	return s.ticketRepo.UpdateFields(ctx, ticketID, updates)
}

// ============================================
// MULTI-LEVEL SUPPORT IMPLEMENTATION
// ============================================

// AutoAssignTicket assigns ticket to best available agent in department
func (s *ticketService) AutoAssignTicket(ctx context.Context, ticketID string, departmentID string) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Get best available agent
	agent, err := s.GetBestAvailableAgent(ctx, departmentID, ticket.CurrentLevel)
	if err != nil {
		return fmt.Errorf("no available agent: %w", err)
	}

	// Assign ticket
	updates := map[string]interface{}{
		"assigned_to_id": agent.ID,
		"department_id":  departmentID,
		"status":         "in_progress",
	}

	if err := s.ticketRepo.UpdateFields(ctx, ticketID, updates); err != nil {
		return fmt.Errorf("failed to assign ticket: %w", err)
	}

	// Create history record
	description := fmt.Sprintf("Auto-assigned to %s (%s)", agent.Name, agent.Email)
	return s.CreateHistoryRecord(ctx, ticketID, "assigned", description, "system")
}

// GetBestAvailableAgent finds the best agent based on workload and level
func (s *ticketService) GetBestAvailableAgent(ctx context.Context, departmentID string, supportLevel int) (*domain.User, error) {
	// Get all agents in department
	agents, err := s.userRepo.GetAgentsByDepartment(ctx, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available in department")
	}

	// Convert string to UUID for department lookup
	deptUUID, err := uuid.Parse(departmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Get department to check assignment rule
	dept, err := s.departmentRepo.GetByID(ctx, deptUUID)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	switch dept.AutoAssignmentRule {
	case "least_loaded":
		return s.getLeastLoadedAgent(ctx, agents)
	case "skill_based":
		return s.getSkillBasedAgent(ctx, agents, supportLevel)
	default: // round_robin
		return s.getRoundRobinAgent(ctx, agents)
	}
}

// ValidateEscalation checks if escalation is valid
func (s *ticketService) ValidateEscalation(ctx context.Context, ticketID string, toLevel int) error {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if !ticket.CanEscalate {
		return fmt.Errorf("ticket cannot be escalated")
	}

	if toLevel <= ticket.CurrentLevel {
		return fmt.Errorf("can only escalate to higher level")
	}

	if toLevel > ticket.MaxLevel {
		return fmt.Errorf("cannot escalate beyond max level %d", ticket.MaxLevel)
	}

	return nil
}

// GetEscalationPath finds the target department for escalation
func (s *ticketService) GetEscalationPath(ctx context.Context, fromDeptID string, toLevel int) (*domain.Department, error) {
	// Convert string to UUID
	fromDeptUUID, err := uuid.Parse(fromDeptID)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Get current department
	currentDept, err := s.departmentRepo.GetByID(ctx, fromDeptUUID)
	if err != nil {
		return nil, fmt.Errorf("current department not found: %w", err)
	}

	// If has explicit escalation department, use it
	if currentDept.EscalationDeptID.Valid {
		targetDeptUUID, err := uuid.Parse(currentDept.EscalationDeptID.String)
		if err != nil {
			return nil, fmt.Errorf("invalid escalation department ID: %w", err)
		}

		targetDept, err := s.departmentRepo.GetByID(ctx, targetDeptUUID)
		if err != nil {
			return nil, fmt.Errorf("escalation department not found: %w", err)
		}

		// Validate target department can handle the level
		if targetDept.SupportLevel == toLevel {
			return targetDept, nil
		}
	}

	// Find department with matching support level
	departments, err := s.departmentRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments: %w", err)
	}

	for _, dept := range departments {
		if dept.SupportLevel == toLevel && dept.CanHandleTickets && dept.IsActive {
			return dept, nil
		}
	}

	return nil, fmt.Errorf("no department available for level L%d", toLevel)
}

// CreateHistoryRecord creates audit trail entry
func (s *ticketService) CreateHistoryRecord(ctx context.Context, ticketID, action, description, createdBy string) error {
	history := &domain.TicketHistory{
		TicketID:    ticketID,
		Action:      action,
		Description: description,
		ActorName:   createdBy,
		UserID:      sql.NullString{String: createdBy, Valid: createdBy != "system"},
	}

	return s.historyRepo.Create(ctx, history)
}

// CheckSLABreach checks if ticket has breached SLA
func (s *ticketService) CheckSLABreach(ctx context.Context, ticketID string) (bool, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return false, fmt.Errorf("ticket not found: %w", err)
	}

	if ticket.SLA == nil {
		return false, nil // No SLA defined
	}

	now := time.Now()

	// Check first response SLA
	if !ticket.FirstResponseAt.Valid {
		if ticket.SLA.FirstResponseDue.Before(now) {
			return true, nil
		}
	}

	// Check resolution SLA
	if !ticket.SLA.ResolvedAt.Valid {
		if ticket.SLA.ResolutionDue.Before(now) {
			return true, nil
		}
	}

	return false, nil
}

// TriggerAutoEscalation automatically escalates ticket
func (s *ticketService) TriggerAutoEscalation(ctx context.Context, ticketID string, reason string) error {
	// Auto-escalate to next level
	req := &domain.EscalateTicketRequest{
		TicketID:    ticketID,
		Reason:      fmt.Sprintf("Auto-escalation: %s", reason),
		EscalatedBy: "system",
	}

	return s.EscalateTicket(ctx, req)
}

// ============================================
// HELPER METHODS
// ============================================

// getDepName safely gets department name
func getDepName(dept *domain.Department) string {
	if dept != nil {
		return dept.Name
	}
	return "Unknown"
}

// adjustSLAForEscalation adjusts SLA timelines after escalation
func (s *ticketService) adjustSLAForEscalation(ctx context.Context, ticket *domain.Ticket, newLevel int) error {
	if ticket.SLA == nil {
		return nil
	}

	// Get category to check if SLA should be adjusted
	if ticket.Category != nil {
		// Higher levels may get more time for resolution
		escalationMultiplier := 1.0 + (float64(newLevel) * 0.5) // 50% more time per level

		newResolutionDue := time.Now().Add(
			time.Duration(float64(ticket.Category.SLAResolution)*escalationMultiplier) * time.Second,
		)

		return s.slaRepo.Update(ctx, &domain.TicketSLA{
			ID:            ticket.SLA.ID,
			ResolutionDue: newResolutionDue,
		})
	}

	return nil
}

// Assignment algorithm implementations
func (s *ticketService) getLeastLoadedAgent(ctx context.Context, agents []*domain.User) (*domain.User, error) {
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// For now, return first agent (can be enhanced with actual workload checking)
	// TODO: Implement workload counting from tickets table
	return agents[0], nil
}

func (s *ticketService) getSkillBasedAgent(ctx context.Context, agents []*domain.User, supportLevel int) (*domain.User, error) {
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// For now, return first agent (can be enhanced with skill matching)
	// TODO: Implement skill-based matching based on agent profile and support level
	return agents[0], nil
}

func (s *ticketService) getRoundRobinAgent(ctx context.Context, agents []*domain.User) (*domain.User, error) {
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// For now, return first agent (can be enhanced with round-robin state)
	// TODO: Implement proper round-robin with state tracking
	return agents[0], nil
}

// ============================================
// STATISTICS & ANALYTICS IMPLEMENTATION
// ============================================

// GetDashboardStats returns comprehensive dashboard statistics
func (s *ticketService) GetDashboardStats(ctx context.Context) (*domain.TicketDashboardStats, error) {
	// For now, use simple counts from existing methods
	// This is a basic implementation - can be enhanced later

	// Use the existing GetList method to get counts
	allFilter := &domain.TicketFilter{Limit: 1, Offset: 0}
	_, totalTickets, err := s.ticketRepo.GetList(ctx, allFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count total tickets: %w", err)
	}

	// Get open tickets
	openFilter := &domain.TicketFilter{Status: &[]string{"open"}[0], Limit: 1, Offset: 0}
	_, openTickets, err := s.ticketRepo.GetList(ctx, openFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count open tickets: %w", err)
	}

	// Get in-progress tickets
	inProgressStatus := "in_progress"
	inProgressFilter := &domain.TicketFilter{Status: &inProgressStatus, Limit: 1, Offset: 0}
	_, inProgressTickets, err := s.ticketRepo.GetList(ctx, inProgressFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count in-progress tickets: %w", err)
	}

	// Get resolved tickets
	resolvedStatus := "resolved"
	resolvedFilter := &domain.TicketFilter{Status: &resolvedStatus, Limit: 1, Offset: 0}
	_, resolvedTickets, err := s.ticketRepo.GetList(ctx, resolvedFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count resolved tickets: %w", err)
	}

	// Get closed tickets
	closedStatus := "closed"
	closedFilter := &domain.TicketFilter{Status: &closedStatus, Limit: 1, Offset: 0}
	_, closedTickets, err := s.ticketRepo.GetList(ctx, closedFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count closed tickets: %w", err)
	}

	// Get today's created tickets
	today := time.Now().Truncate(24 * time.Hour)
	todayFilter := &domain.TicketFilter{DateFrom: &today, Limit: 1, Offset: 0}
	_, todayCreated, err := s.ticketRepo.GetList(ctx, todayFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count today's tickets: %w", err)
	}

	// Basic implementation - will be enhanced with proper statistics methods
	return &domain.TicketDashboardStats{
		TotalTickets:         int(totalTickets),
		OpenTickets:          int(openTickets),
		InProgressTickets:    int(inProgressTickets),
		ResolvedTickets:      int(resolvedTickets),
		ClosedTickets:        int(closedTickets),
		EscalatedTickets:     0,   // Will be implemented later
		OverdueTickets:       0,   // Will be implemented later
		AvgResolutionTime:    0.0, // Will be implemented later
		AvgFirstResponseTime: 0.0, // Will be implemented later
		SLAComplianceRate:    0.0, // Will be implemented later
		TodayCreated:         int(todayCreated),
		TodayResolved:        0, // Will be implemented later
		TicketsByStatus:      make(map[string]int),
		TicketsByPriority:    make(map[string]int),
		TicketsByCategory:    make(map[string]int),
		TicketsByLevel:       make(map[string]int),
	}, nil
}

// GetTicketStats returns detailed ticket statistics
func (s *ticketService) GetTicketStats(ctx context.Context, req *domain.TicketStatsRequest) (*domain.DetailedTicketStats, error) {
	// Implementation placeholder - this would involve complex queries
	// For now, return basic structure
	return &domain.DetailedTicketStats{
		TotalCount:        0,
		StatusBreakdown:   make(map[string]int),
		PriorityBreakdown: make(map[string]int),
		CategoryBreakdown: make(map[string]int),
		LevelBreakdown:    make(map[string]int),
		DailyTrends:       []domain.TicketDailyTrend{},
		ResolutionMetrics: domain.TicketResolutionMetrics{},
		SLAMetrics:        domain.TicketSLAMetrics{},
		EscalationMetrics: domain.TicketEscalationMetrics{},
	}, nil
}

// GetPerformanceMetrics returns agent and department performance metrics
func (s *ticketService) GetPerformanceMetrics(ctx context.Context, req *domain.TicketPerformanceRequest) (*domain.TicketPerformanceMetrics, error) {
	// Implementation placeholder - this would involve complex performance calculations
	return &domain.TicketPerformanceMetrics{
		AgentPerformance:      []domain.TicketAgentPerformance{},
		DepartmentPerformance: []domain.TicketDepartmentPerformance{},
		OverallMetrics:        domain.TicketOverallPerformance{},
	}, nil
}

// GetSLAReport returns SLA compliance reports
func (s *ticketService) GetSLAReport(ctx context.Context, req *domain.TicketSLAReportRequest) (*domain.TicketSLAReport, error) {
	// Implementation placeholder - this would involve SLA analysis
	return &domain.TicketSLAReport{
		Period:          req.Period,
		TotalTickets:    0,
		TicketsWithSLA:  0,
		ComplianceStats: domain.TicketSLACompliance{},
		BreachAnalysis:  domain.TicketBreachAnalysis{},
		TrendData:       []domain.TicketSLATrend{},
	}, nil
}

// GetEscalationStats returns escalation statistics and trends
func (s *ticketService) GetEscalationStats(ctx context.Context, req *domain.TicketEscalationStatsRequest) (*domain.TicketEscalationStats, error) {
	// Implementation placeholder - this would involve escalation analysis
	return &domain.TicketEscalationStats{
		Period:             req.Period,
		TotalEscalations:   0,
		EscalationRate:     0.0,
		LevelBreakdown:     make(map[string]int),
		TriggerBreakdown:   make(map[string]int),
		SuccessRate:        0.0,
		AvgEscalationTime:  0.0,
		TrendData:          []domain.TicketEscalationTrendData{},
		TopEscalationPaths: []domain.TicketEscalationPath{},
	}, nil
}

// ============================================
// BULK OPERATIONS IMPLEMENTATION
// ============================================

// BulkAssignTickets assigns multiple tickets to an agent
func (s *ticketService) BulkAssignTickets(ctx context.Context, req *domain.BulkAssignRequest) error {
	// Validate agent exists
	_, err := s.userRepo.GetByID(ctx, req.AgentID)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Process each ticket
	for _, ticketID := range req.TicketIDs {
		assignReq := &domain.AssignTicketRequest{
			TicketID:   ticketID,
			AgentID:    req.AgentID,
			AssignedBy: req.AssignedBy,
		}

		if err := s.AssignTicket(ctx, assignReq); err != nil {
			// Log error but continue with other tickets
			continue
		}

		// Create history record
		description := fmt.Sprintf("Bulk assigned - %s", req.Reason)
		s.CreateHistoryRecord(ctx, ticketID, "bulk_assigned", description, req.AssignedBy)
	}

	return nil
}

// BulkUpdateStatus updates status for multiple tickets
func (s *ticketService) BulkUpdateStatus(ctx context.Context, req *domain.BulkUpdateStatusRequest) error {
	// Process each ticket
	for _, ticketID := range req.TicketIDs {
		updateReq := &domain.UpdateTicketRequest{
			ID:     ticketID,
			Status: &req.Status,
		}

		if _, err := s.UpdateTicket(ctx, updateReq); err != nil {
			// Log error but continue with other tickets
			continue
		}

		// Create history record
		description := fmt.Sprintf("Bulk status update to %s - %s", req.Status, req.Reason)
		s.CreateHistoryRecord(ctx, ticketID, "bulk_status_update", description, req.UpdatedBy)
	}

	return nil
}

// BulkCloseTickets closes multiple tickets
func (s *ticketService) BulkCloseTickets(ctx context.Context, req *domain.BulkCloseRequest) error {
	// Process each ticket
	for _, ticketID := range req.TicketIDs {
		status := "closed"
		updateReq := &domain.UpdateTicketRequest{
			ID:     ticketID,
			Status: &status,
		}

		if _, err := s.UpdateTicket(ctx, updateReq); err != nil {
			// Log error but continue with other tickets
			continue
		}

		// Create history record
		description := fmt.Sprintf("Bulk closed - %s", req.CloseReason)
		s.CreateHistoryRecord(ctx, ticketID, "bulk_closed", description, req.ClosedBy)
	}

	return nil
}
