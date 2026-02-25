package ports

import (
	"context"
	"nexus-super-node-v3/internal/core/domain"
)

// AgentRepository defines the contract for persisting agent data
type AgentRepository interface {
	Create(ctx context.Context, agent *domain.Agent) error
	GetByID(ctx context.Context, id string) (*domain.Agent, error)
	List(ctx context.Context, ownerID string) ([]*domain.Agent, error)
	Update(ctx context.Context, agent *domain.Agent) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status domain.AgentStatus) error
}

// AgentService defines the business logic for managing agents
type AgentService interface {
	CreateAgent(ctx context.Context, agent *domain.Agent) error
	GetAgent(ctx context.Context, id string) (*domain.Agent, error)
	ListAgents(ctx context.Context, ownerID string) ([]*domain.Agent, error)
	UpdateAgent(ctx context.Context, agent *domain.Agent) error
	DeleteAgent(ctx context.Context, id string) error
	DeployAgent(ctx context.Context, id string) error
	PauseAgent(ctx context.Context, id string) error
}
