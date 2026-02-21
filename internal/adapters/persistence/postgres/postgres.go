package postgres

import (
	"context"
	"database/sql"
	"nexus-super-node-v3/internal/adapters/persistence/postgres/db"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// postgresRepository is the implementation of the UserRepository port.
type postgresRepository struct {
	db db.Querier
}

// NewPostgresRepository creates a new Postgres repository.
func NewPostgresRepository(ctx context.Context, dsn string) (ports.UserRepository, error) {
	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &postgresRepository{
		db: db.New(dbConn),
	}, nil
}

// GetUser retrieves a user by ID.
func (r *postgresRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	user, err := r.db.GetUser(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}

	return toDomainUser(user), nil
}

// ListUsers retrieves all users.
func (r *postgresRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := r.db.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	domainUsers := make([]*domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = toDomainUser(user)
	}

	return domainUsers, nil
}

// CreateUser creates a new user.
func (r *postgresRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	createdUser, err := r.db.CreateUser(ctx, user.Address)
	if err != nil {
		return nil, err
	}

	return toDomainUser(createdUser), nil
}

func (r *postgresRepository) GetUserByAddress(ctx context.Context, address string) (*domain.User, error) {
	user, err := r.db.GetUserByAddress(ctx, address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}

	return toDomainUser(user), nil
}

// toDomainUser converts a db.User to a domain.User.
func toDomainUser(user db.User) *domain.User {
	return &domain.User{
		ID:      int64(user.ID),
		Address: user.Address,
	}
}
