package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Filter structs for repository queries
type TicketFilter struct {
	Status       *string
	Priority     *string
	CategoryID   *string
	DepartmentID *string
	AssignedTo   *string
	CreatedBy    *string
	CustomerInfo *string // Search in customer_name, customer_email, customer_phone
	Subject      *string // Search in subject
	DateFrom     *time.Time
	DateTo       *time.Time
	Limit        int
	Offset       int
	SortBy       string // "created_at", "updated_at", "priority", "status"
	SortOrder    string // "asc", "desc"
}

// Request/Response structs for ticket service
type CreateTicketRequest struct {
	Subject       string  `json:"subject" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	CustomerName  string  `json:"customer_name" validate:"required"`
	CustomerEmail string  `json:"customer_email" validate:"required,email"`
	CustomerPhone string  `json:"customer_phone"`
	Priority      string  `json:"priority" validate:"required,oneof=low medium high urgent"`
	CategoryID    string  `json:"category_id" validate:"required"`
	CreatedBy     *string `json:"created_by"` // Agent ID if created by agent
	CreatedVia    string  `json:"created_via" validate:"required,oneof=customer agent ai"`
	AutoAssign    bool    `json:"auto_assign"`
}

type UpdateTicketRequest struct {
	ID          string  `json:"id" validate:"required"`
	Subject     *string `json:"subject"`
	Description *string `json:"description"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	Status      *string `json:"status" validate:"omitempty,oneof=open in_progress resolved closed escalated"`
	AssignedTo  *string `json:"assigned_to"`
	UpdatedBy   *string `json:"updated_by"`
}

type AssignTicketRequest struct {
	TicketID   string `json:"ticket_id" validate:"required"`
	AgentID    string `json:"agent_id" validate:"required"`
	AssignedBy string `json:"assigned_by" validate:"required"`
}

type EscalateTicketRequest struct {
	TicketID    string  `json:"ticket_id" validate:"required"`
	ToLevel     *int    `json:"to_level"` // Optional: specify target level
	Reason      string  `json:"reason" validate:"required"`
	EscalatedBy string  `json:"escalated_by" validate:"required"`
	AssignToID  *string `json:"assign_to_id"`        // Optional - specify agent, otherwise auto-assign
	ForceDept   *string `json:"force_department_id"` // Optional: force to specific dept
}

type AddCommentRequest struct {
	TicketID        string `json:"ticket_id" validate:"required"`
	Content         string `json:"content" validate:"required"`
	IsPublic        bool   `json:"is_public"`
	CreatedBy       string `json:"created_by" validate:"required"`
	IsAgentResponse bool   `json:"is_agent_response"`
}

// TicketEscalation Filter and Stats structs
type TicketEscalationFilter struct {
	TicketID         string
	FromLevel        *int
	ToLevel          *int
	FromDepartmentID string
	ToDepartmentID   string
	EscalatedByID    string
	IsAutoEscalation *bool
	TriggerType      string
	WasSuccessful    *bool
	StartDate        time.Time
	EndDate          time.Time
	Limit            int
	Offset           int
}

type EscalationLevelStats struct {
	Level                 int     `json:"level"`
	TotalEscalations      int64   `json:"total_escalations"`
	SuccessfulEscalations int64   `json:"successful_escalations"`
	AutoEscalations       int64   `json:"auto_escalations"`
	AvgResolutionHours    float64 `json:"avg_resolution_hours"`
}

type EscalationTrend struct {
	Date                 time.Time `json:"date"`
	TotalEscalations     int64     `json:"total_escalations"`
	SLABreachEscalations int64     `json:"sla_breach_escalations"`
	ManualEscalations    int64     `json:"manual_escalations"`
	AutoEscalations      int64     `json:"auto_escalations"`
}

type EscalationReason struct {
	TriggerType string  `json:"trigger_type"`
	Count       int64   `json:"count"`
	Percentage  float64 `json:"percentage"`
}

// Usecase interfaces
type OSSChatUsecase interface {
	StartChat(ctx context.Context, req *StartChatRequest, ipAddress string) (*StartChatResponse, error)
	SetSessionContact(ctx context.Context, req *SetSessionContactRequest) (*SetSessionContactResponse, error)
	LinkOSSUser(ctx context.Context, req *LinkOSSUserRequest) (*LinkOSSUserResponse, error)
	GetChatHistory(ctx context.Context, req *GetChatHistoryRequest) (*GetChatHistoryResponse, error)
}

// Repository interfaces
// UserRepository interface for user operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	GetAgentsByDepartment(ctx context.Context, departmentID string) ([]*User, error)
	GetAvailableAgents(ctx context.Context, departmentID *string) ([]*User, error)
	GetWithPagination(ctx context.Context, offset, limit int, role string, departmentID *string) ([]*User, error)
	Count(ctx context.Context, role string, departmentID *string) (int, error)
	GetByRole(ctx context.Context, role string) ([]*User, error)
	// Analytics methods
	CountByRole(ctx context.Context, role string) (int64, error)
	CountOnlineAgents(ctx context.Context) (int64, error)
}

// DepartmentRepository interface for department operations
type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	GetAll(ctx context.Context) ([]*Department, error)
	GetWithPagination(ctx context.Context, offset, limit int, search string, isActive *bool, supportLevel *int, parentDeptID string) ([]*Department, error)
	GetWithRelations(ctx context.Context, id uuid.UUID) (*Department, error)
	GetByParent(ctx context.Context, parentDeptID uuid.UUID) ([]*Department, error)
	GetBySupportLevel(ctx context.Context, supportLevel int) ([]*Department, error)
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, search string, isActive *bool, supportLevel *int, parentDeptID string) (int64, error)
}

// ChatUserRepository interface for chat user operations
type ChatUserRepository interface {
	Create(ctx context.Context, user *ChatUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatUser, error)
	GetByBrowserUUID(ctx context.Context, browserUUID uuid.UUID) (*ChatUser, error)
	GetByOSSUserID(ctx context.Context, ossUserID string) (*ChatUser, error)
	GetByEmail(ctx context.Context, email string) (*ChatUser, error)
	Update(ctx context.Context, user *ChatUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	LinkOSSUser(ctx context.Context, browserUUID uuid.UUID, ossUserID string, email string) error
	List(ctx context.Context, limit, offset int) ([]*ChatUser, error)
	Count(ctx context.Context) (int64, error)
}

// ChatSessionContactRepository interface for chat session contact operations
type ChatSessionContactRepository interface {
	Create(ctx context.Context, contact *ChatSessionContact) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*ChatSessionContact, error)
	Update(ctx context.Context, contact *ChatSessionContact) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
}

// ChatSessionRepository interface for chat session operations
type ChatSessionRepository interface {
	Create(ctx context.Context, session *ChatSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatSession, error)
	GetByChatUserID(ctx context.Context, chatUserID uuid.UUID) ([]*ChatSession, error)
	GetByAgentID(ctx context.Context, agentID uuid.UUID) ([]*ChatSession, error)
	GetActiveSessions(ctx context.Context) ([]*ChatSession, error)
	GetWaitingSessions(ctx context.Context) ([]*ChatSession, error)
	Update(ctx context.Context, session *ChatSession) error
	Close(ctx context.Context, sessionID uuid.UUID) error
	GetSessionsByStatus(ctx context.Context, status string) ([]*ChatSession, error)
	GetSessionsByDateRange(ctx context.Context, start, end time.Time) ([]*ChatSession, error)
	GetWithPagination(ctx context.Context, offset, limit int, status string, agentID, departmentID *uuid.UUID) ([]*ChatSession, error)
	GetSessionsWithMessages(ctx context.Context, chatUserID uuid.UUID, limit, offset int) ([]*ChatSession, error)
	GetSessionHistory(ctx context.Context, chatUserID uuid.UUID, limit, offset int) ([]*ChatSession, error)
	Count(ctx context.Context, status string, agentID, departmentID *uuid.UUID) (int, error)
	// Analytics methods
	CountByStatus(ctx context.Context, status string) (int64, error)
	CountCompletedSince(ctx context.Context, since time.Time) (int64, error)
	GetAverageResponseTime(ctx context.Context) (float64, error)
	GetOSSCategoriesStats(ctx context.Context) ([]CategoryStats, error)
}

// ChatMessageRepository interface for chat message operations
type ChatMessageRepository interface {
	Create(ctx context.Context, message *ChatMessage) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatMessage, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*ChatMessage, error)
	Update(ctx context.Context, message *ChatMessage) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMessagesByDateRange(ctx context.Context, start, end time.Time) ([]*ChatMessage, error)
	// Analytics methods
	GetTopQuestions(ctx context.Context, limit int) ([]QuestionStats, error)
}

// ChatLogRepository interface for chat log operations
type ChatLogRepository interface {
	Create(ctx context.Context, log *ChatLog) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*ChatLog, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*ChatLog, error)
}

// ChatTagRepository interface for chat tag operations
type ChatTagRepository interface {
	Create(ctx context.Context, tag *ChatTag) error
	GetByID(ctx context.Context, id uuid.UUID) (*ChatTag, error)
	GetAll(ctx context.Context) ([]*ChatTag, error)
	Update(ctx context.Context, tag *ChatTag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ChatSessionTagRepository interface for chat session tag operations
type ChatSessionTagRepository interface {
	Create(ctx context.Context, sessionTag *ChatSessionTag) error
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*ChatSessionTag, error)
	DeleteBySessionID(ctx context.Context, sessionID uuid.UUID) error
	DeleteByTagID(ctx context.Context, tagID uuid.UUID) error
}

// AgentStatusRepository interface for agent status operations
type AgentStatusRepository interface {
	Create(ctx context.Context, status *AgentStatus) error
	GetByAgentID(ctx context.Context, agentID uuid.UUID) (*AgentStatus, error)
	Update(ctx context.Context, status *AgentStatus) error
	GetOnlineAgents(ctx context.Context) ([]*AgentStatus, error)
	UpdateLastActive(ctx context.Context, agentID uuid.UUID) error
}

// ChatAnalyticsRepository interface for chat analytics operations
type ChatAnalyticsRepository interface {
	Create(ctx context.Context, analytics *ChatAnalytics) error
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*ChatAnalytics, error)
	GetByAgentAndDateRange(ctx context.Context, agentID uuid.UUID, start, end time.Time) ([]*ChatAnalytics, error)
	GetByDepartmentAndDateRange(ctx context.Context, departmentID uuid.UUID, start, end time.Time) ([]*ChatAnalytics, error)
	UpdateOrCreate(ctx context.Context, analytics *ChatAnalytics) error
}

// ============================================
// TICKETING SYSTEM REPOSITORIES
// ============================================

// TicketRepository interface for ticket operations
type TicketRepository interface {
	Create(ctx context.Context, ticket *Ticket) error
	GetByID(ctx context.Context, id string) (*Ticket, error)
	GetByTicketCode(ctx context.Context, code string) (*Ticket, error)
	GetByAccessToken(ctx context.Context, token string) (*Ticket, error)
	Update(ctx context.Context, ticket *Ticket) error
	UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetList(ctx context.Context, filter *TicketFilter) ([]*Ticket, int64, error)
	GetByAssignedAgent(ctx context.Context, agentID string, status []string) ([]*Ticket, error)
	GetByDepartment(ctx context.Context, deptID string, status []string) ([]*Ticket, error)
	GetOverdueTickets(ctx context.Context) ([]*Ticket, error)

	// Statistics Methods
	CountByStatus(ctx context.Context, statuses []string) (int64, error)
	CountEscalatedTickets(ctx context.Context) (int64, error)
	CountOverdueTickets(ctx context.Context) (int64, error)
	CountCreatedSince(ctx context.Context, since time.Time) (int64, error)
	CountResolvedSince(ctx context.Context, since time.Time) (int64, error)
	GetAverageResolutionTime(ctx context.Context) (float64, error)
	GetAverageFirstResponseTime(ctx context.Context) (float64, error)
	GetSLAComplianceRate(ctx context.Context) (float64, error)
	GetStatusBreakdown(ctx context.Context) (map[string]int, error)
	GetPriorityBreakdown(ctx context.Context) (map[string]int, error)
	GetCategoryBreakdown(ctx context.Context) (map[string]int, error)
	GetLevelBreakdown(ctx context.Context) (map[string]int, error)
}

// TicketCategoryRepository interface for ticket category operations
type TicketCategoryRepository interface {
	Create(ctx context.Context, category *TicketCategory) error
	GetByID(ctx context.Context, id string) (*TicketCategory, error)
	GetByCode(ctx context.Context, code string) (*TicketCategory, error)
	GetAll(ctx context.Context) ([]*TicketCategory, error)
	GetActive(ctx context.Context) ([]*TicketCategory, error)
	Update(ctx context.Context, category *TicketCategory) error
	Delete(ctx context.Context, id string) error
}

// TicketCommentRepository interface for ticket comment operations
type TicketCommentRepository interface {
	Create(ctx context.Context, comment *TicketComment) error
	GetByID(ctx context.Context, id string) (*TicketComment, error)
	GetByTicketID(ctx context.Context, ticketID string, includeInternal bool) ([]*TicketComment, error)
	Update(ctx context.Context, comment *TicketComment) error
	Delete(ctx context.Context, id string) error
}

// TicketAttachmentRepository interface for ticket attachment operations
type TicketAttachmentRepository interface {
	Create(ctx context.Context, attachment *TicketAttachment) error
	GetByID(ctx context.Context, id string) (*TicketAttachment, error)
	GetByTicketID(ctx context.Context, ticketID string) ([]*TicketAttachment, error)
	Delete(ctx context.Context, id string) error
}

// TicketHistoryRepository interface for ticket history operations
type TicketHistoryRepository interface {
	Create(ctx context.Context, history *TicketHistory) error
	GetByTicketID(ctx context.Context, ticketID string) ([]*TicketHistory, error)
}

// TicketSLARepository interface for ticket SLA operations
type TicketSLARepository interface {
	Create(ctx context.Context, sla *TicketSLA) error
	GetByTicketID(ctx context.Context, ticketID string) (*TicketSLA, error)
	Update(ctx context.Context, sla *TicketSLA) error
	GetBreachedTickets(ctx context.Context) ([]*TicketSLA, error)
}

// TicketEscalationRepository interface for ticket escalation operations
type TicketEscalationRepository interface {
	Create(ctx context.Context, escalation *TicketEscalation) error
	GetByID(ctx context.Context, id string) (*TicketEscalation, error)
	GetByTicketID(ctx context.Context, ticketID string) ([]TicketEscalation, error)
	GetAll(ctx context.Context, filter *TicketEscalationFilter) ([]TicketEscalation, error)
	Update(ctx context.Context, escalation *TicketEscalation) error
	UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *TicketEscalationFilter) (int64, error)
	GetEscalationStatsByLevel(ctx context.Context, departmentID string, days int) ([]EscalationLevelStats, error)
	GetEscalationTrends(ctx context.Context, days int) ([]EscalationTrend, error)
	GetTopEscalationReasons(ctx context.Context, limit int, days int) ([]EscalationReason, error)
}

// EmailService interface for email operations
type EmailService interface {
	SendEmail(ctx context.Context, request *SendEmailRequest) (*EmailResponse, error)
	SendTemplatedEmail(ctx context.Context, template *EmailTemplate, to []string, variables map[string]string) (*EmailResponse, error)
	SendWelcomeEmail(ctx context.Context, to string, name string) (*EmailResponse, error)
	SendPasswordResetEmail(ctx context.Context, to string, resetToken string) (*EmailResponse, error)
	SendChatTranscriptEmail(ctx context.Context, to string, transcript string, sessionID uuid.UUID) (*EmailResponse, error)
}

// TicketService interface for ticket business logic
type TicketService interface {
	CreateTicket(ctx context.Context, req *CreateTicketRequest) (*Ticket, error)
	GetTicket(ctx context.Context, id string) (*Ticket, error)
	GetTicketByCode(ctx context.Context, code string) (*Ticket, error)
	GetTicketByAccessToken(ctx context.Context, token string) (*Ticket, error)
	UpdateTicket(ctx context.Context, req *UpdateTicketRequest) (*Ticket, error)
	AssignTicket(ctx context.Context, req *AssignTicketRequest) error
	EscalateTicket(ctx context.Context, req *EscalateTicketRequest) error
	AddComment(ctx context.Context, req *AddCommentRequest) (*TicketComment, error)
	GetTicketList(ctx context.Context, filter *TicketFilter) ([]*Ticket, int64, error)
	GetAgentTickets(ctx context.Context, agentID string, status []string) ([]*Ticket, error)
	GetDepartmentTickets(ctx context.Context, deptID string, status []string) ([]*Ticket, error)

	// Multi-Level Support Methods
	AutoAssignTicket(ctx context.Context, ticketID string, departmentID string) error
	GetBestAvailableAgent(ctx context.Context, departmentID string, supportLevel int) (*User, error)
	ValidateEscalation(ctx context.Context, ticketID string, toLevel int) error
	GetEscalationPath(ctx context.Context, fromDeptID string, toLevel int) (*Department, error)
	CreateHistoryRecord(ctx context.Context, ticketID, action, description, createdBy string) error
	CheckSLABreach(ctx context.Context, ticketID string) (bool, error)
	TriggerAutoEscalation(ctx context.Context, ticketID string, reason string) error

	// Statistics & Analytics Methods
	GetDashboardStats(ctx context.Context) (*TicketDashboardStats, error)
	GetTicketStats(ctx context.Context, req *TicketStatsRequest) (*DetailedTicketStats, error)
	GetPerformanceMetrics(ctx context.Context, req *TicketPerformanceRequest) (*TicketPerformanceMetrics, error)
	GetSLAReport(ctx context.Context, req *TicketSLAReportRequest) (*TicketSLAReport, error)
	GetEscalationStats(ctx context.Context, req *TicketEscalationStatsRequest) (*TicketEscalationStats, error)

	// Bulk Operations
	BulkAssignTickets(ctx context.Context, req *BulkAssignRequest) error
	BulkUpdateStatus(ctx context.Context, req *BulkUpdateStatusRequest) error
	BulkCloseTickets(ctx context.Context, req *BulkCloseRequest) error
}
