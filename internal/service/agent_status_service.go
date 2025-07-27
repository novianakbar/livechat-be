package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/internal/infrastructure/repository"
)

type AgentStatusService struct {
	agentStatusRepo *repository.AgentStatusRepository
	userRepo        domain.UserRepository
}

func NewAgentStatusService(
	agentStatusRepo *repository.AgentStatusRepository,
	userRepo domain.UserRepository,
) *AgentStatusService {
	return &AgentStatusService{
		agentStatusRepo: agentStatusRepo,
		userRepo:        userRepo,
	}
}

// UpdateAgentHeartbeat updates agent's heartbeat and status
func (s *AgentStatusService) UpdateAgentHeartbeat(ctx context.Context, agentID uuid.UUID, status string) error {
	// Get agent details from database
	agent, err := s.userRepo.GetByID(ctx, agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent details: %w", err)
	}

	if agent == nil {
		return fmt.Errorf("agent not found")
	}

	// Verify user is an agent
	if agent.Role != "agent" && agent.Role != "admin" {
		return fmt.Errorf("user is not an agent")
	}

	// Validate status
	validStatuses := map[string]bool{
		"online": true,
		"busy":   true,
		"away":   true,
	}

	if !validStatuses[status] {
		status = "online" // Default to online
	}

	// Update status in Redis
	return s.agentStatusRepo.SetAgentOnline(ctx, agent, status)
}

// GetOnlineAgents gets all currently online agents
func (s *AgentStatusService) GetOnlineAgents(ctx context.Context) ([]repository.AgentOnlineStatus, error) {
	return s.agentStatusRepo.GetAllOnlineAgents(ctx)
}

// GetOnlineAgentsByDepartment gets online agents by department
func (s *AgentStatusService) GetOnlineAgentsByDepartment(ctx context.Context, departmentID uuid.UUID) ([]repository.AgentOnlineStatus, error) {
	return s.agentStatusRepo.GetOnlineAgentsByDepartment(ctx, departmentID)
}

// GetAgentStatus gets specific agent status
func (s *AgentStatusService) GetAgentStatus(ctx context.Context, agentID uuid.UUID) (*repository.AgentOnlineStatus, error) {
	return s.agentStatusRepo.GetAgentStatus(ctx, agentID)
}

// SetAgentOffline removes agent from online status
func (s *AgentStatusService) SetAgentOffline(ctx context.Context, agentID uuid.UUID) error {
	return s.agentStatusRepo.RemoveAgentOnline(ctx, agentID)
}

// GetDepartmentStats gets online agents count by department
func (s *AgentStatusService) GetDepartmentStats(ctx context.Context) (map[string]int, error) {
	return s.agentStatusRepo.GetDepartmentStats(ctx)
}
