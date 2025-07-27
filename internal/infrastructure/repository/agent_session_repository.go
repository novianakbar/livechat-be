package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-shared/entities"
	"gorm.io/gorm"
)

type AgentSessionRepository struct {
	db *gorm.DB
}

func NewAgentSessionRepository(db *gorm.DB) *AgentSessionRepository {
	return &AgentSessionRepository{
		db: db,
	}
}

// SetAgentLoggedIn records agent login in database
func (r *AgentSessionRepository) SetAgentLoggedIn(ctx context.Context, agentID uuid.UUID) error {
	agentStatus := &entities.AgentStatus{
		AgentID: agentID,
		Status:  "logged_in",
	}

	// Use UPSERT to handle existing records
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Assign("status = ?", "logged_in").
		FirstOrCreate(agentStatus).Error

	if err != nil {
		return fmt.Errorf("failed to set agent logged in: %w", err)
	}

	return nil
}

// SetAgentLoggedOut records agent logout in database
func (r *AgentSessionRepository) SetAgentLoggedOut(ctx context.Context, agentID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&entities.AgentStatus{}).
		Where("agent_id = ?", agentID).
		Updates(map[string]interface{}{
			"status":     "logged_out",
			"updated_at": "NOW()",
		}).Error

	if err != nil {
		return fmt.Errorf("failed to set agent logged out: %w", err)
	}

	return nil
}

// GetAgentSessionStatus gets agent login session status from database
func (r *AgentSessionRepository) GetAgentSessionStatus(ctx context.Context, agentID uuid.UUID) (*entities.AgentStatus, error) {
	var agentStatus entities.AgentStatus
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		First(&agentStatus).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get agent session status: %w", err)
	}

	return &agentStatus, nil
}

// GetLoggedInAgents gets all currently logged in agents
func (r *AgentSessionRepository) GetLoggedInAgents(ctx context.Context) ([]entities.AgentStatus, error) {
	var agentStatuses []entities.AgentStatus
	err := r.db.WithContext(ctx).
		Where("status = ?", "logged_in").
		Preload("Agent").
		Find(&agentStatuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get logged in agents: %w", err)
	}

	return agentStatuses, nil
}

// GetLoggedInAgentsByDepartment gets logged in agents by department
func (r *AgentSessionRepository) GetLoggedInAgentsByDepartment(ctx context.Context, departmentID uuid.UUID) ([]entities.AgentStatus, error) {
	var agentStatuses []entities.AgentStatus
	err := r.db.WithContext(ctx).
		Joins("JOIN users ON agent_status.agent_id = users.id").
		Where("agent_status.status = ? AND users.department_id = ?", "logged_in", departmentID).
		Preload("Agent").
		Find(&agentStatuses).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get logged in agents by department: %w", err)
	}

	return agentStatuses, nil
}
