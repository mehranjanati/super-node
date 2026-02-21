
package gateway

import (
	"context"
	"nexus-super-node-v3/internal/core/domain"
)

// MockUserRepository is a mock implementation of the UserRepository interface.

type MockUserRepository struct {
	GetUserFunc        func(ctx context.Context, id int64) (*domain.User, error)
	ListUsersFunc      func(ctx context.Context) ([]*domain.User, error)
	CreateUserFunc     func(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByAddressFunc func(ctx context.Context, address string) (*domain.User, error)
}

func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return m.GetUserFunc(ctx, id)
}

func (m *MockUserRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return m.ListUsersFunc(ctx)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return m.CreateUserFunc(ctx, user)
}

func (m *MockUserRepository) GetUserByAddress(ctx context.Context, address string) (*domain.User, error) {
	return m.GetUserByAddressFunc(ctx, address)
}
