package ports

import (
	"context"

	"nexus-super-node-v3/internal/core/domain"
)

// UserRepository is the port for interacting with user data.

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByAddress(ctx context.Context, address string) (*domain.User, error)
}

// AppDataRepository is the port for interacting with app data.

type AppDataRepository interface {
	GetAppData(ctx context.Context, id string) (*domain.AppData, error)
	CreateAppData(ctx context.Context, id string, data []byte) error
}

// Database is the port for interacting with the database.

type Database interface {
	UserRepository
	AppDataRepository
}
