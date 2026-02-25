package agent

import (
	"context"
	"errors"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"
	"time"

	"github.com/google/uuid"
)

type agentService struct {
	repo ports.AgentRepository
}

func NewAgentService(repo ports.AgentRepository) ports.AgentService {
	return &agentService{
		repo: repo,
	}
}

func (s *agentService) CreateAgent(ctx context.Context, agent *domain.Agent) error {
	if agent.ID == "" {
		agent.ID = uuid.New().String()
	}
	if agent.Name == "" {
		return errors.New("agent name is required")
	}
	if agent.Status == "" {
		agent.Status = domain.AgentStatusDeploying
	}
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()

	// Initialize default performance metrics
	agent.Performance = domain.AgentPerformance{
		ROI:         0,
		Trades:      0,
		Uptime:      100,
		SuccessRate: 0,
		LastActive:  time.Now(),
	}

	return s.repo.Create(ctx, agent)
}

func (s *agentService) GetAgent(ctx context.Context, id string) (*domain.Agent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *agentService) ListAgents(ctx context.Context, ownerID string) ([]*domain.Agent, error) {
	return s.repo.List(ctx, ownerID)
}

func (s *agentService) UpdateAgent(ctx context.Context, agent *domain.Agent) error {
	agent.UpdatedAt = time.Now()
	return s.repo.Update(ctx, agent)
}

func (s *agentService) DeleteAgent(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *agentService) DeployAgent(ctx context.Context, id string) error {
	// TODO: Trigger actual deployment (Temporal workflow, etc.)
	return s.repo.UpdateStatus(ctx, id, domain.AgentStatusDeploying)
}

func (s *agentService) PauseAgent(ctx context.Context, id string) error {
	// TODO: Stop running instances
	return s.repo.UpdateStatus(ctx, id, domain.AgentStatusPaused)
}
