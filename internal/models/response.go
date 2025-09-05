package models

// ChatUserResponse represents a clean chat user response
type ChatUserResponse struct {
	ID          string `json:"id"`
	BrowserUUID string `json:"browser_uuid,omitempty"`
	OSSUserID   string `json:"oss_user_id,omitempty"`
	Email       string `json:"email,omitempty"`
	IsAnonymous bool   `json:"is_anonymous"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UserResponse represents a clean user response
type UserResponse struct {
	ID           string              `json:"id"`
	Email        string              `json:"email"`
	Name         string              `json:"name"`
	Role         string              `json:"role"`
	IsActive     bool                `json:"is_active"`
	DepartmentID string              `json:"department_id,omitempty"`
	Department   *DepartmentResponse `json:"department,omitempty"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

// DepartmentResponse represents a clean department response
type DepartmentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`

	// Ticketing specific fields
	CanHandleTickets   bool `json:"can_handle_tickets"`
	MaxTicketsPerAgent int  `json:"max_tickets_per_agent"`

	// Multi-Level Support Fields
	SupportLevel       int                 `json:"support_level"`
	ParentDeptID       string              `json:"parent_dept_id,omitempty"`
	ParentDept         *DepartmentResponse `json:"parent_dept,omitempty"`
	MaxEscalationLevel int                 `json:"max_escalation_level"`
	AutoAssignmentRule string              `json:"auto_assignment_rule"`
	EscalationDeptID   string              `json:"escalation_dept_id,omitempty"`
	EscalationDept     *DepartmentResponse `json:"escalation_dept,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ChatSessionContactResponse represents a clean session contact response
type ChatSessionContactResponse struct {
	ID           string `json:"id"`
	SessionID    string `json:"session_id"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Position     string `json:"position,omitempty"`
	CompanyName  string `json:"company_name,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ChatMessageResponse represents a clean chat message response
type ChatMessageResponse struct {
	ID          string   `json:"id"`
	SessionID   string   `json:"session_id"`
	SenderID    string   `json:"sender_id,omitempty"`
	SenderType  string   `json:"sender_type"`
	Message     string   `json:"message"`
	MessageType string   `json:"message_type"`
	Attachments []string `json:"attachments,omitempty"`
	ReadAt      string   `json:"read_at,omitempty"`
	CreatedAt   string   `json:"created_at"`
}

// ChatSessionMinimalResponse represents a minimal chat session response (for lists)
type ChatSessionMinimalResponse struct {
	ID         string            `json:"id"`
	ChatUserID string            `json:"chat_user_id"`
	AgentID    string            `json:"agent_id,omitempty"`
	Topic      string            `json:"topic"`
	Status     string            `json:"status"`
	Priority   string            `json:"priority"`
	StartedAt  string            `json:"started_at"`
	EndedAt    string            `json:"ended_at,omitempty"`
	ChatUser   *ChatUserResponse `json:"chat_user,omitempty"`
	Agent      *UserResponse     `json:"agent,omitempty"`
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
}

// ChatSessionDetailResponse represents a detailed chat session response (for single session)
type ChatSessionDetailResponse struct {
	ID           string                      `json:"id"`
	ChatUserID   string                      `json:"chat_user_id"`
	AgentID      string                      `json:"agent_id,omitempty"`
	DepartmentID string                      `json:"department_id,omitempty"`
	Topic        string                      `json:"topic"`
	Status       string                      `json:"status"`
	Priority     string                      `json:"priority"`
	StartedAt    string                      `json:"started_at"`
	EndedAt      string                      `json:"ended_at,omitempty"`
	ChatUser     *ChatUserResponse           `json:"chat_user,omitempty"`
	Agent        *UserResponse               `json:"agent,omitempty"`
	Department   *DepartmentResponse         `json:"department,omitempty"`
	Messages     []ChatMessageResponse       `json:"messages,omitempty"`
	Contact      *ChatSessionContactResponse `json:"contact,omitempty"`
	CreatedAt    string                      `json:"created_at"`
	UpdatedAt    string                      `json:"updated_at"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// TicketResponse represents a clean ticket response
type TicketResponse struct {
	ID              string                  `json:"id"`
	TicketCode      string                  `json:"ticket_code"`
	Subject         string                  `json:"subject"`
	Description     string                  `json:"description"`
	CustomerName    string                  `json:"customer_name"`
	CustomerEmail   string                  `json:"customer_email"`
	CustomerPhone   string                  `json:"customer_phone,omitempty"`
	Priority        string                  `json:"priority"`
	Status          string                  `json:"status"`
	CreatedVia      string                  `json:"created_via"`
	CategoryID      string                  `json:"category_id,omitempty"`
	Category        *TicketCategoryResponse `json:"category,omitempty"`
	AssignedToID    string                  `json:"assigned_to_id,omitempty"`
	AssignedTo      *UserResponse           `json:"assigned_to,omitempty"`
	DepartmentID    string                  `json:"department_id,omitempty"`
	Department      *DepartmentResponse     `json:"department,omitempty"`
	CreatedByID     string                  `json:"created_by_id,omitempty"`
	CreatedBy       *UserResponse           `json:"created_by,omitempty"`
	AccessToken     string                  `json:"access_token,omitempty"`
	FirstResponseAt string                  `json:"first_response_at,omitempty"`
	ResolvedAt      string                  `json:"resolved_at,omitempty"`
	ClosedAt        string                  `json:"closed_at,omitempty"`
	CreatedAt       string                  `json:"created_at"`
	UpdatedAt       string                  `json:"updated_at"`
}

// TicketCategoryResponse represents a clean ticket category response
type TicketCategoryResponse struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Code                string              `json:"code"`
	Description         string              `json:"description,omitempty"`
	Color               string              `json:"color,omitempty"`
	SLAFirstResponse    int                 `json:"sla_first_response"`
	SLAResolution       int                 `json:"sla_resolution"`
	DefaultDepartmentID string              `json:"default_department_id,omitempty"`
	DefaultDepartment   *DepartmentResponse `json:"default_department,omitempty"`
	IsActive            bool                `json:"is_active"`
	CreatedAt           string              `json:"created_at"`
	UpdatedAt           string              `json:"updated_at"`
}

// TicketCommentResponse represents a clean ticket comment response
type TicketCommentResponse struct {
	ID             string        `json:"id"`
	TicketID       string        `json:"ticket_id"`
	UserID         string        `json:"user_id,omitempty"`
	User           *UserResponse `json:"user,omitempty"`
	Content        string        `json:"content"`
	IsInternal     bool          `json:"is_internal"`
	IsFromCustomer bool          `json:"is_from_customer"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
}

// TicketAttachmentResponse represents a clean ticket attachment response
type TicketAttachmentResponse struct {
	ID          string `json:"id"`
	TicketID    string `json:"ticket_id"`
	FileName    string `json:"file_name"`
	FilePath    string `json:"file_path"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
	UploadedBy  string `json:"uploaded_by,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// TicketHistoryResponse represents a clean ticket history response
type TicketHistoryResponse struct {
	ID          string        `json:"id"`
	TicketID    string        `json:"ticket_id"`
	UserID      string        `json:"user_id,omitempty"`
	User        *UserResponse `json:"user,omitempty"`
	Action      string        `json:"action"`
	OldValue    string        `json:"old_value,omitempty"`
	NewValue    string        `json:"new_value,omitempty"`
	Description string        `json:"description,omitempty"`
	CreatedAt   string        `json:"created_at"`
}

// TicketSLAResponse represents a clean ticket SLA response
type TicketSLAResponse struct {
	ID                      string `json:"id"`
	TicketID                string `json:"ticket_id"`
	FirstResponseDue        string `json:"first_response_due"`
	FirstResponseAt         string `json:"first_response_at,omitempty"`
	ResolutionDue           string `json:"resolution_due"`
	ResolutionAt            string `json:"resolution_at,omitempty"`
	IsFirstResponseBreached bool   `json:"is_first_response_breached"`
	IsResolutionBreached    bool   `json:"is_resolution_breached"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
}

// TicketDetailResponse represents a detailed ticket response with all relations
type TicketDetailResponse struct {
	*TicketResponse
	Comments    []TicketCommentResponse    `json:"comments,omitempty"`
	Attachments []TicketAttachmentResponse `json:"attachments,omitempty"`
	History     []TicketHistoryResponse    `json:"history,omitempty"`
	SLA         *TicketSLAResponse         `json:"sla,omitempty"`
}
