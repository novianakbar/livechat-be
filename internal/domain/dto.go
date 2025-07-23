package domain

import (
	"time"

	"github.com/google/uuid"
)

// Auth DTOs
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *User     `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RegisterRequest struct {
	Email        string     `json:"email" validate:"required,email"`
	Password     string     `json:"password" validate:"required,min=6"`
	Name         string     `json:"name" validate:"required"`
	Role         string     `json:"role" validate:"required,oneof=admin agent"`
	DepartmentID *uuid.UUID `json:"department_id"`
}

// Department DTOs
type CreateDepartmentRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateDepartmentRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// Chat Session DTOs
type StartChatRequest struct {
	CompanyName string `json:"company_name" validate:"required"`
	PersonName  string `json:"person_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Topic       string `json:"topic" validate:"required"`
	Priority    string `json:"priority" validate:"oneof=low normal high urgent"`
}

type StartChatResponse struct {
	SessionID uuid.UUID `json:"session_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
}

type AssignAgentRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	AgentID   uuid.UUID `json:"agent_id" validate:"required"`
}

type TransferChatRequest struct {
	SessionID       uuid.UUID  `json:"session_id" validate:"required"`
	NewAgentID      *uuid.UUID `json:"new_agent_id"`
	NewDepartmentID *uuid.UUID `json:"new_department_id"`
	Reason          string     `json:"reason"`
}

// Chat Message DTOs
type SendMessageRequest struct {
	SessionID   uuid.UUID `json:"session_id" validate:"required"`
	Message     string    `json:"message" validate:"required"`
	MessageType string    `json:"message_type" validate:"oneof=text image file system"`
	Attachments []string  `json:"attachments"`
}

type SendMessageResponse struct {
	MessageID uuid.UUID `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// Tag DTOs
type CreateTagRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color"`
}

type AddTagToSessionRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	TagID     uuid.UUID `json:"tag_id" validate:"required"`
}

// Analytics DTOs
type AnalyticsRequest struct {
	StartDate    time.Time  `json:"start_date" validate:"required"`
	EndDate      time.Time  `json:"end_date" validate:"required"`
	AgentID      *uuid.UUID `json:"agent_id"`
	DepartmentID *uuid.UUID `json:"department_id"`
}

type AnalyticsResponse struct {
	TotalSessions       int                   `json:"total_sessions"`
	CompletedSessions   int                   `json:"completed_sessions"`
	AverageResponseTime float64               `json:"average_response_time"`
	TotalMessages       int                   `json:"total_messages"`
	SessionsByStatus    map[string]int        `json:"sessions_by_status"`
	SessionsByPriority  map[string]int        `json:"sessions_by_priority"`
	DailyAnalytics      []DailyAnalytics      `json:"daily_analytics"`
	AgentPerformance    []AgentPerformance    `json:"agent_performance"`
	DepartmentAnalytics []DepartmentAnalytics `json:"department_analytics"`
}

type DailyAnalytics struct {
	Date              time.Time `json:"date"`
	TotalSessions     int       `json:"total_sessions"`
	CompletedSessions int       `json:"completed_sessions"`
	TotalMessages     int       `json:"total_messages"`
}

type AgentPerformance struct {
	Agent               *User   `json:"agent"`
	TotalSessions       int     `json:"total_sessions"`
	CompletedSessions   int     `json:"completed_sessions"`
	AverageResponseTime float64 `json:"average_response_time"`
	TotalMessages       int     `json:"total_messages"`
}

type DepartmentAnalytics struct {
	Department        *Department `json:"department"`
	TotalSessions     int         `json:"total_sessions"`
	CompletedSessions int         `json:"completed_sessions"`
	TotalMessages     int         `json:"total_messages"`
	ActiveAgents      int         `json:"active_agents"`
}

// Agent Status DTOs
type UpdateAgentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=online offline busy away"`
}

type AgentStatusResponse struct {
	Agent        *User     `json:"agent"`
	Status       string    `json:"status"`
	LastActiveAt time.Time `json:"last_active_at"`
}

// WebSocket DTOs
type WebSocketMessage struct {
	Type      string      `json:"type"`
	SessionID uuid.UUID   `json:"session_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type WebSocketResponse struct {
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

// Common DTOs
type PaginationRequest struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Search   string `json:"search" query:"search"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDir  string `json:"sort_dir" query:"sort_dir"`
}

type PaginationResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Session Management DTOs
type SessionListRequest struct {
	Status       string     `json:"status" query:"status"`
	AgentID      *uuid.UUID `json:"agent_id" query:"agent_id"`
	DepartmentID *uuid.UUID `json:"department_id" query:"department_id"`
	CustomerID   *uuid.UUID `json:"customer_id" query:"customer_id"`
	Priority     string     `json:"priority" query:"priority"`
	DateFrom     *time.Time `json:"date_from" query:"date_from"`
	DateTo       *time.Time `json:"date_to" query:"date_to"`
	PaginationRequest
}

type CloseSessionRequest struct {
	SessionID uuid.UUID `json:"session_id" validate:"required"`
	Reason    string    `json:"reason"`
	Rating    *int      `json:"rating" validate:"omitempty,min=1,max=5"`
	Feedback  string    `json:"feedback"`
}

// Analytics related structures for OSS support system
type DashboardStats struct {
	ActiveSessions      int             `json:"activeSessions"`
	WaitingSessions     int             `json:"waitingSessions"`
	CompletedToday      int             `json:"completedToday"`
	AverageResponseTime int             `json:"averageResponseTime"` // in seconds
	TotalAgents         int             `json:"totalAgents"`
	OnlineAgents        int             `json:"onlineAgents"`
	TopQuestions        []QuestionStats `json:"topQuestions"`
	OSSCategories       []CategoryStats `json:"ossCategories"`
}

type QuestionStats struct {
	Question string `json:"question"`
	Count    int    `json:"count"`
}

type CategoryStats struct {
	Category   string `json:"category"`
	Count      int    `json:"count"`
	Percentage int    `json:"percentage"`
}

type GetAnalyticsRequest struct {
	StartDate    *string `query:"start_date"`
	EndDate      *string `query:"end_date"`
	DepartmentID *string `query:"department_id"`
	AgentID      *string `query:"agent_id"`
}

// Pagination related DTOs
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message,omitempty"`
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// Email DTOs
type EmailRequest struct {
	To      []string `json:"to" validate:"required,min=1,dive,email"`
	Subject string   `json:"subject" validate:"required"`
	Content string   `json:"content" validate:"required"`
	IsHTML  bool     `json:"is_html"`
}

type SendEmailRequest struct {
	To        []string          `json:"to" validate:"required,min=1,dive,email"`
	Subject   string            `json:"subject" validate:"required"`
	Content   string            `json:"content" validate:"required"`
	IsHTML    bool              `json:"is_html"`
	Variables map[string]string `json:"variables,omitempty"`
}

type EmailTemplate struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	IsHTML  bool   `json:"is_html"`
}

type EmailResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Message   string `json:"message,omitempty"`
}
