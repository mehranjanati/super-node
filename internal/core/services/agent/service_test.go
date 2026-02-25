package agent

import (
	"context"
	"errors"
	"testing"

	"nexus-super-node-v3/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAgentRepository is a mock implementation of ports.AgentRepository
type MockAgentRepository struct {
	mock.Mock
}

func (m *MockAgentRepository) Create(ctx context.Context, agent *domain.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockAgentRepository) GetByID(ctx context.Context, id string) (*domain.Agent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Agent), args.Error(1)
}

func (m *MockAgentRepository) Update(ctx context.Context, agent *domain.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockAgentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAgentRepository) List(ctx context.Context, ownerID string) ([]*domain.Agent, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Agent), args.Error(1)
}

func (m *MockAgentRepository) UpdateStatus(ctx context.Context, id string, status domain.AgentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestCreateAgent(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAgentService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		agent := &domain.Agent{
			Name: "Test Agent",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Agent")).Return(nil)

		err := service.CreateAgent(ctx, agent)

		assert.NoError(t, err)
		assert.NotEmpty(t, agent.ID)
		assert.Equal(t, domain.AgentStatusDeploying, agent.Status)
		assert.NotZero(t, agent.CreatedAt)
		assert.NotZero(t, agent.UpdatedAt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing Name", func(t *testing.T) {
		agent := &domain.Agent{
			Name: "",
		}

		err := service.CreateAgent(ctx, agent)

		assert.Error(t, err)
		assert.Equal(t, "agent name is required", err.Error())
	})
}

func TestGetAgent(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAgentService(mockRepo)
	ctx := context.Background()

	t.Run("Found", func(t *testing.T) {
		expectedAgent := &domain.Agent{
			ID:   "test-id",
			Name: "Test Agent",
		}

		mockRepo.On("GetByID", ctx, "test-id").Return(expectedAgent, nil)

		agent, err := service.GetAgent(ctx, "test-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedAgent, agent)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "unknown-id").Return(nil, errors.New("not found"))

		agent, err := service.GetAgent(ctx, "unknown-id")

		assert.Error(t, err)
		assert.Nil(t, agent)
		mockRepo.AssertExpectations(t)
	})
}

func TestListAgents(t *testing.T) {
	mockRepo := new(MockAgentRepository)
	service := NewAgentService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedAgents := []*domain.Agent{
			{ID: "1", Name: "Agent 1"},
			{ID: "2", Name: "Agent 2"},
		}

		mockRepo.On("List", ctx, "owner-1").Return(expectedAgents, nil)

		agents, err := service.ListAgents(ctx, "owner-1")

		assert.NoError(t, err)
		assert.Equal(t, expectedAgents, agents)
		mockRepo.AssertExpectations(t)
	})
}
