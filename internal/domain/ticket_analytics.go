package domain

import "time"

// ============================================
// TICKETING STATISTICS & ANALYTICS TYPES
// ============================================

// Ticketing Dashboard Statistics
type TicketDashboardStats struct {
	TotalTickets         int            `json:"total_tickets"`
	OpenTickets          int            `json:"open_tickets"`
	InProgressTickets    int            `json:"in_progress_tickets"`
	ResolvedTickets      int            `json:"resolved_tickets"`
	ClosedTickets        int            `json:"closed_tickets"`
	EscalatedTickets     int            `json:"escalated_tickets"`
	OverdueTickets       int            `json:"overdue_tickets"`
	AvgResolutionTime    float64        `json:"avg_resolution_time"`     // hours
	AvgFirstResponseTime float64        `json:"avg_first_response_time"` // hours
	SLAComplianceRate    float64        `json:"sla_compliance_rate"`     // percentage
	TodayCreated         int            `json:"today_created"`
	TodayResolved        int            `json:"today_resolved"`
	TicketsByStatus      map[string]int `json:"tickets_by_status"`
	TicketsByPriority    map[string]int `json:"tickets_by_priority"`
	TicketsByCategory    map[string]int `json:"tickets_by_category"`
	TicketsByLevel       map[string]int `json:"tickets_by_level"`
}

// Statistics Request
type TicketStatsRequest struct {
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	DepartmentID string     `json:"department_id,omitempty"`
	AgentID      string     `json:"agent_id,omitempty"`
}

// Detailed Ticket Statistics
type DetailedTicketStats struct {
	TotalCount        int                     `json:"total_count"`
	StatusBreakdown   map[string]int          `json:"status_breakdown"`
	PriorityBreakdown map[string]int          `json:"priority_breakdown"`
	CategoryBreakdown map[string]int          `json:"category_breakdown"`
	LevelBreakdown    map[string]int          `json:"level_breakdown"`
	DailyTrends       []TicketDailyTrend      `json:"daily_trends"`
	ResolutionMetrics TicketResolutionMetrics `json:"resolution_metrics"`
	SLAMetrics        TicketSLAMetrics        `json:"sla_metrics"`
	EscalationMetrics TicketEscalationMetrics `json:"escalation_metrics"`
}

// Daily Ticket Trend
type TicketDailyTrend struct {
	Date     time.Time `json:"date"`
	Created  int       `json:"created"`
	Resolved int       `json:"resolved"`
	Closed   int       `json:"closed"`
}

// Resolution Metrics
type TicketResolutionMetrics struct {
	AvgResolutionTime    float64 `json:"avg_resolution_time"`    // hours
	MedianResolutionTime float64 `json:"median_resolution_time"` // hours
	MinResolutionTime    float64 `json:"min_resolution_time"`    // hours
	MaxResolutionTime    float64 `json:"max_resolution_time"`    // hours
	ResolutionRate       float64 `json:"resolution_rate"`        // percentage
}

// SLA Metrics
type TicketSLAMetrics struct {
	TotalWithSLA          int     `json:"total_with_sla"`
	FirstResponseOnTime   int     `json:"first_response_on_time"`
	ResolutionOnTime      int     `json:"resolution_on_time"`
	FirstResponseBreached int     `json:"first_response_breached"`
	ResolutionBreached    int     `json:"resolution_breached"`
	ComplianceRate        float64 `json:"compliance_rate"`         // percentage
	AvgFirstResponseTime  float64 `json:"avg_first_response_time"` // hours
}

// Escalation Metrics
type TicketEscalationMetrics struct {
	TotalEscalations      int            `json:"total_escalations"`
	EscalationRate        float64        `json:"escalation_rate"`     // percentage
	AvgEscalationTime     float64        `json:"avg_escalation_time"` // hours
	EscalationsByLevel    map[string]int `json:"escalations_by_level"`
	EscalationsByReason   map[string]int `json:"escalations_by_reason"`
	SuccessfulEscalations int            `json:"successful_escalations"`
}

// Performance Request
type TicketPerformanceRequest struct {
	Period       int    `json:"period"` // days
	DepartmentID string `json:"department_id,omitempty"`
}

// Performance Metrics
type TicketPerformanceMetrics struct {
	AgentPerformance      []TicketAgentPerformance      `json:"agent_performance"`
	DepartmentPerformance []TicketDepartmentPerformance `json:"department_performance"`
	OverallMetrics        TicketOverallPerformance      `json:"overall_metrics"`
}

// Agent Performance
type TicketAgentPerformance struct {
	AgentID              string  `json:"agent_id"`
	AgentName            string  `json:"agent_name"`
	TicketsAssigned      int     `json:"tickets_assigned"`
	TicketsResolved      int     `json:"tickets_resolved"`
	TicketsClosed        int     `json:"tickets_closed"`
	AvgResolutionTime    float64 `json:"avg_resolution_time"`     // hours
	AvgFirstResponseTime float64 `json:"avg_first_response_time"` // hours
	SLAComplianceRate    float64 `json:"sla_compliance_rate"`     // percentage
	CurrentWorkload      int     `json:"current_workload"`
	PerformanceScore     float64 `json:"performance_score"` // 0-100
}

// Department Performance
type TicketDepartmentPerformance struct {
	DepartmentID      string  `json:"department_id"`
	DepartmentName    string  `json:"department_name"`
	TicketsHandled    int     `json:"tickets_handled"`
	TicketsResolved   int     `json:"tickets_resolved"`
	TicketsEscalated  int     `json:"tickets_escalated"`
	AvgResolutionTime float64 `json:"avg_resolution_time"` // hours
	SLAComplianceRate float64 `json:"sla_compliance_rate"` // percentage
	EscalationRate    float64 `json:"escalation_rate"`     // percentage
	CurrentWorkload   int     `json:"current_workload"`
	PerformanceScore  float64 `json:"performance_score"` // 0-100
}

// Overall Performance
type TicketOverallPerformance struct {
	TotalTicketsHandled   int     `json:"total_tickets_handled"`
	AvgResolutionTime     float64 `json:"avg_resolution_time"`     // hours
	OverallSLACompliance  float64 `json:"overall_sla_compliance"`  // percentage
	OverallEscalationRate float64 `json:"overall_escalation_rate"` // percentage
	TopPerformingAgent    string  `json:"top_performing_agent"`
	TopPerformingDept     string  `json:"top_performing_dept"`
}

// SLA Report Request
type TicketSLAReportRequest struct {
	Period       int    `json:"period"` // days
	DepartmentID string `json:"department_id,omitempty"`
}

// SLA Report
type TicketSLAReport struct {
	Period          int                  `json:"period"`
	TotalTickets    int                  `json:"total_tickets"`
	TicketsWithSLA  int                  `json:"tickets_with_sla"`
	ComplianceStats TicketSLACompliance  `json:"compliance_stats"`
	BreachAnalysis  TicketBreachAnalysis `json:"breach_analysis"`
	TrendData       []TicketSLATrend     `json:"trend_data"`
}

// SLA Compliance
type TicketSLACompliance struct {
	FirstResponseCompliance float64 `json:"first_response_compliance"` // percentage
	ResolutionCompliance    float64 `json:"resolution_compliance"`     // percentage
	OverallCompliance       float64 `json:"overall_compliance"`        // percentage
}

// Breach Analysis
type TicketBreachAnalysis struct {
	FirstResponseBreaches int                  `json:"first_response_breaches"`
	ResolutionBreaches    int                  `json:"resolution_breaches"`
	TotalBreaches         int                  `json:"total_breaches"`
	BreachRate            float64              `json:"breach_rate"` // percentage
	TopBreachReasons      []TicketBreachReason `json:"top_breach_reasons"`
}

// Breach Reason
type TicketBreachReason struct {
	Reason string  `json:"reason"`
	Count  int     `json:"count"`
	Rate   float64 `json:"rate"` // percentage
}

// SLA Trend
type TicketSLATrend struct {
	Date         time.Time `json:"date"`
	Compliance   float64   `json:"compliance"` // percentage
	Breaches     int       `json:"breaches"`
	TotalTickets int       `json:"total_tickets"`
}

// Escalation Stats Request
type TicketEscalationStatsRequest struct {
	Period       int    `json:"period"` // days
	DepartmentID string `json:"department_id,omitempty"`
}

// Escalation Stats
type TicketEscalationStats struct {
	Period             int                         `json:"period"`
	TotalEscalations   int                         `json:"total_escalations"`
	EscalationRate     float64                     `json:"escalation_rate"` // percentage
	LevelBreakdown     map[string]int              `json:"level_breakdown"`
	TriggerBreakdown   map[string]int              `json:"trigger_breakdown"`
	SuccessRate        float64                     `json:"success_rate"`        // percentage
	AvgEscalationTime  float64                     `json:"avg_escalation_time"` // hours
	TrendData          []TicketEscalationTrendData `json:"trend_data"`
	TopEscalationPaths []TicketEscalationPath      `json:"top_escalation_paths"`
}

// Escalation Trend Data
type TicketEscalationTrendData struct {
	Date        time.Time `json:"date"`
	Escalations int       `json:"escalations"`
	SuccessRate float64   `json:"success_rate"` // percentage
}

// Escalation Path
type TicketEscalationPath struct {
	FromLevel   int     `json:"from_level"`
	ToLevel     int     `json:"to_level"`
	Count       int     `json:"count"`
	SuccessRate float64 `json:"success_rate"` // percentage
}

// ============================================
// BULK OPERATIONS TYPES
// ============================================

// Bulk Assign Request
type BulkAssignRequest struct {
	TicketIDs    []string `json:"ticket_ids" validate:"required,min=1"`
	AgentID      string   `json:"agent_id" validate:"required"`
	DepartmentID string   `json:"department_id,omitempty"`
	Reason       string   `json:"reason,omitempty"`
	AssignedBy   string   `json:"assigned_by" validate:"required"`
}

// Bulk Update Status Request
type BulkUpdateStatusRequest struct {
	TicketIDs []string `json:"ticket_ids" validate:"required,min=1"`
	Status    string   `json:"status" validate:"required,oneof=open in_progress resolved closed"`
	Reason    string   `json:"reason,omitempty"`
	UpdatedBy string   `json:"updated_by" validate:"required"`
}

// Bulk Close Request
type BulkCloseRequest struct {
	TicketIDs   []string `json:"ticket_ids" validate:"required,min=1"`
	CloseReason string   `json:"close_reason,omitempty"`
	Resolution  string   `json:"resolution,omitempty"`
	ClosedBy    string   `json:"closed_by" validate:"required"`
}
