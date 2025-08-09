package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/redis/go-redis/v9"
)

type AgentStatusRepository struct {
	redisClient *redis.Client
}

func NewAgentStatusRepository(redisClient *redis.Client) *AgentStatusRepository {
	return &AgentStatusRepository{
		redisClient: redisClient,
	}
}

// AgentOnlineStatus represents agent status stored in Redis
type AgentOnlineStatus struct {
	AgentID       string    `json:"agent_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	DepartmentID  *string   `json:"department_id"`
	Department    string    `json:"department"`
	Status        string    `json:"status"` // online, busy, away
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

const (
	agentOnlinePrefix  = "agent:online:"
	agentsByDeptPrefix = "agents:dept:"
	allAgentsKey       = "agents:all"
	agentStatusTTL     = 5 * time.Minute // Agent considered offline after 5 minutes
)

// SetAgentOnline sets agent status as online with heartbeat
func (r *AgentStatusRepository) SetAgentOnline(ctx context.Context, agent *domain.User, status string) error {
	var departmentID *string
	if agent.DepartmentID.Valid {
		departmentID = &agent.DepartmentID.String
	}

	agentStatus := AgentOnlineStatus{
		AgentID:       agent.ID,
		Name:          agent.Name,
		Email:         agent.Email,
		DepartmentID:  departmentID,
		Status:        status,
		LastHeartbeat: time.Now(),
	}

	// Add department name if available
	if agent.Department != nil {
		agentStatus.Department = agent.Department.Name
	}

	statusJSON, err := json.Marshal(agentStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal agent status: %w", err)
	}

	pipe := r.redisClient.Pipeline()

	// Set individual agent status with TTL
	agentKey := fmt.Sprintf("%s%s", agentOnlinePrefix, agent.ID)
	pipe.Set(ctx, agentKey, statusJSON, agentStatusTTL)

	// Add to all agents set
	pipe.SAdd(ctx, allAgentsKey, agent.ID)
	pipe.Expire(ctx, allAgentsKey, agentStatusTTL*2) // Keep the set longer

	// Add to department-specific set if agent has department
	if agent.DepartmentID.Valid {
		deptKey := fmt.Sprintf("%s%s", agentsByDeptPrefix, agent.DepartmentID.String)
		pipe.SAdd(ctx, deptKey, agent.ID)
		pipe.Expire(ctx, deptKey, agentStatusTTL*2)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// GetAgentStatus gets specific agent status
func (r *AgentStatusRepository) GetAgentStatus(ctx context.Context, agentID uuid.UUID) (*AgentOnlineStatus, error) {
	agentKey := fmt.Sprintf("%s%s", agentOnlinePrefix, agentID.String())
	statusJSON, err := r.redisClient.Get(ctx, agentKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Agent not online
		}
		return nil, fmt.Errorf("failed to get agent status: %w", err)
	}

	var status AgentOnlineStatus
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal agent status: %w", err)
	}

	return &status, nil
}

// GetAllOnlineAgents gets all currently online agents
func (r *AgentStatusRepository) GetAllOnlineAgents(ctx context.Context) ([]AgentOnlineStatus, error) {
	agentIDs, err := r.redisClient.SMembers(ctx, allAgentsKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []AgentOnlineStatus{}, nil
		}
		return nil, fmt.Errorf("failed to get agent IDs: %w", err)
	}

	if len(agentIDs) == 0 {
		return []AgentOnlineStatus{}, nil
	}

	// Get all agent statuses
	keys := make([]string, len(agentIDs))
	for i, agentID := range agentIDs {
		keys[i] = fmt.Sprintf("%s%s", agentOnlinePrefix, agentID)
	}

	statusStrings, err := r.redisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get agent statuses: %w", err)
	}

	var agents []AgentOnlineStatus
	for i, statusStr := range statusStrings {
		if statusStr == nil {
			// Agent no longer online, remove from set
			r.redisClient.SRem(ctx, allAgentsKey, agentIDs[i])
			continue
		}

		var status AgentOnlineStatus
		if err := json.Unmarshal([]byte(statusStr.(string)), &status); err != nil {
			continue // Skip invalid entries
		}

		agents = append(agents, status)
	}

	return agents, nil
}

// GetOnlineAgentsByDepartment gets online agents by department
func (r *AgentStatusRepository) GetOnlineAgentsByDepartment(ctx context.Context, departmentID uuid.UUID) ([]AgentOnlineStatus, error) {
	deptKey := fmt.Sprintf("%s%s", agentsByDeptPrefix, departmentID.String())
	agentIDs, err := r.redisClient.SMembers(ctx, deptKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []AgentOnlineStatus{}, nil
		}
		return nil, fmt.Errorf("failed to get department agent IDs: %w", err)
	}

	if len(agentIDs) == 0 {
		return []AgentOnlineStatus{}, nil
	}

	// Get all agent statuses
	keys := make([]string, len(agentIDs))
	for i, agentID := range agentIDs {
		keys[i] = fmt.Sprintf("%s%s", agentOnlinePrefix, agentID)
	}

	statusStrings, err := r.redisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get agent statuses: %w", err)
	}

	var agents []AgentOnlineStatus
	for i, statusStr := range statusStrings {
		if statusStr == nil {
			// Agent no longer online, remove from department set
			r.redisClient.SRem(ctx, deptKey, agentIDs[i])
			continue
		}

		var status AgentOnlineStatus
		if err := json.Unmarshal([]byte(statusStr.(string)), &status); err != nil {
			continue // Skip invalid entries
		}

		agents = append(agents, status)
	}

	return agents, nil
}

// RemoveAgentOnline removes agent from online status
func (r *AgentStatusRepository) RemoveAgentOnline(ctx context.Context, agentID uuid.UUID) error {
	agentKey := fmt.Sprintf("%s%s", agentOnlinePrefix, agentID.String())

	pipe := r.redisClient.Pipeline()

	// Remove agent status
	pipe.Del(ctx, agentKey)

	// Remove from all agents set
	pipe.SRem(ctx, allAgentsKey, agentID.String())

	// We don't know the department, so we'll let the cleanup happen naturally via TTL
	// or when GetOnlineAgentsByDepartment is called

	_, err := pipe.Exec(ctx)
	return err
}

// GetDepartmentStats gets statistics for departments
func (r *AgentStatusRepository) GetDepartmentStats(ctx context.Context) (map[string]int, error) {
	agents, err := r.GetAllOnlineAgents(ctx)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	for _, agent := range agents {
		dept := "No Department"
		if agent.Department != "" {
			dept = agent.Department
		}
		stats[dept]++
	}

	return stats, nil
}
