package models

// CreateDepartmentRequest represents the request for creating a new department
type CreateDepartmentRequest struct {
	Name               string  `json:"name" validate:"required,min=2,max=100"`
	Description        *string `json:"description" validate:"omitempty,max=500"`
	IsActive           *bool   `json:"is_active" validate:"omitempty"`
	CanHandleTickets   *bool   `json:"can_handle_tickets" validate:"omitempty"`
	MaxTicketsPerAgent *int    `json:"max_tickets_per_agent" validate:"omitempty,min=1,max=100"`

	// Multi-Level Support Fields
	SupportLevel       *int    `json:"support_level" validate:"omitempty,min=0,max=3"`
	ParentDeptID       *string `json:"parent_dept_id" validate:"omitempty,uuid"`
	MaxEscalationLevel *int    `json:"max_escalation_level" validate:"omitempty,min=0,max=3"`
	AutoAssignmentRule *string `json:"auto_assignment_rule" validate:"omitempty,oneof=round_robin least_loaded skill_based"`
	EscalationDeptID   *string `json:"escalation_dept_id" validate:"omitempty,uuid"`
}

// UpdateDepartmentRequest represents the request for updating an existing department
type UpdateDepartmentRequest struct {
	Name               *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description        *string `json:"description" validate:"omitempty,max=500"`
	IsActive           *bool   `json:"is_active" validate:"omitempty"`
	CanHandleTickets   *bool   `json:"can_handle_tickets" validate:"omitempty"`
	MaxTicketsPerAgent *int    `json:"max_tickets_per_agent" validate:"omitempty,min=1,max=100"`

	// Multi-Level Support Fields
	SupportLevel       *int    `json:"support_level" validate:"omitempty,min=0,max=3"`
	ParentDeptID       *string `json:"parent_dept_id" validate:"omitempty,uuid"`
	MaxEscalationLevel *int    `json:"max_escalation_level" validate:"omitempty,min=0,max=3"`
	AutoAssignmentRule *string `json:"auto_assignment_rule" validate:"omitempty,oneof=round_robin least_loaded skill_based"`
	EscalationDeptID   *string `json:"escalation_dept_id" validate:"omitempty,uuid"`
}

// DepartmentQueryRequest represents query parameters for filtering departments
type DepartmentQueryRequest struct {
	Page         int    `query:"page" validate:"omitempty,min=1"`
	Limit        int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search       string `query:"search" validate:"omitempty,max=100"`
	IsActive     *bool  `query:"is_active" validate:"omitempty"`
	SupportLevel *int   `query:"support_level" validate:"omitempty,min=0,max=3"`
	ParentDeptID string `query:"parent_dept_id" validate:"omitempty,uuid"`
}
