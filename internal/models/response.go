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
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
