package tidb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"

	_ "github.com/go-sql-driver/mysql"
)

type tidbRepository struct {
	db *sql.DB
}

// NewTiDBRepository creates a new TiDB/MySQL repository
func NewTiDBRepository(ctx context.Context, dsn string) (ports.UserRepository, ports.FinanceRepository, ports.AppDataRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open tidb connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to ping tidb: %w", err)
	}

	// Ensure tables exist (Simple migration for now)
	if err := migrate(ctx, db); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to migrate tidb schema: %w", err)
	}

	repo := &tidbRepository{db: db}
	return repo, repo, repo, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			address VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS assets (
			id VARCHAR(255) PRIMARY KEY,
			symbol VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL,
			total_supply DOUBLE NOT NULL,
			owner_id VARCHAR(255) NOT NULL,
			metadata JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_balances (
			user_id VARCHAR(255) NOT NULL,
			asset_id VARCHAR(255) NOT NULL,
			balance DOUBLE NOT NULL DEFAULT 0,
			PRIMARY KEY (user_id, asset_id)
		);`,
		`CREATE TABLE IF NOT EXISTS loans (
			id VARCHAR(255) PRIMARY KEY,
			borrower_id VARCHAR(255) NOT NULL,
			collateral_id VARCHAR(255) NOT NULL,
			amount DOUBLE NOT NULL,
			interest_rate DOUBLE NOT NULL,
			status VARCHAR(50) NOT NULL,
			due_date TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS rewards (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			asset_id VARCHAR(255) NOT NULL,
			amount DOUBLE NOT NULL,
			reason VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS app_data (
			id VARCHAR(255) PRIMARY KEY,
			data BLOB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

// --- UserRepository Implementation ---

func (r *tidbRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, address FROM users WHERE id = ?", id).Scan(&user.ID, &user.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *tidbRepository) GetUserByAddress(ctx context.Context, address string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, address FROM users WHERE address = ?", address).Scan(&user.ID, &user.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *tidbRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO users (address) VALUES (?)", user.Address)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (r *tidbRepository) ListUsers(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, address FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Address); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// --- FinanceRepository Implementation ---

func (r *tidbRepository) SaveAsset(ctx context.Context, asset *domain.Asset) error {
	// Metadata is JSON, skipping for simplicity in this snippet or handling via a helper if needed.
	// For now inserting NULL for metadata or stringifying it.
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO assets (id, symbol, name, type, total_supply, owner_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		asset.ID, asset.Symbol, asset.Name, asset.Type, asset.TotalSupply, asset.OwnerID, asset.CreatedAt,
	)
	return err
}

func (r *tidbRepository) GetAsset(ctx context.Context, id string) (*domain.Asset, error) {
	var asset domain.Asset
	// Ignoring metadata scan for brevity
	err := r.db.QueryRowContext(ctx, "SELECT id, symbol, name, type, total_supply, owner_id, created_at FROM assets WHERE id = ?", id).
		Scan(&asset.ID, &asset.Symbol, &asset.Name, &asset.Type, &asset.TotalSupply, &asset.OwnerID, &asset.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &asset, nil
}

func (r *tidbRepository) UpdateUserBalance(ctx context.Context, userID, assetID string, delta float64) error {
	// Upsert logic
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_balances (user_id, asset_id, balance) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE balance = balance + ?`,
		userID, assetID, delta, delta,
	)
	return err
}

func (r *tidbRepository) GetBalance(ctx context.Context, userID, assetID string) (*domain.UserBalance, error) {
	var ub domain.UserBalance
	err := r.db.QueryRowContext(ctx, "SELECT user_id, asset_id, balance FROM user_balances WHERE user_id = ? AND asset_id = ?", userID, assetID).
		Scan(&ub.UserID, &ub.AssetID, &ub.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return zero balance if not found
			return &domain.UserBalance{UserID: userID, AssetID: assetID, Balance: 0}, nil
		}
		return nil, err
	}
	return &ub, nil
}

func (r *tidbRepository) SaveLoan(ctx context.Context, loan *domain.Loan) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO loans (id, borrower_id, collateral_id, amount, interest_rate, status, due_date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		loan.ID, loan.BorrowerID, loan.CollateralID, loan.Amount, loan.InterestRate, loan.Status, loan.DueDate,
	)
	return err
}

func (r *tidbRepository) UpdateLoan(ctx context.Context, loan *domain.Loan) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE loans SET status = ? WHERE id = ?",
		loan.Status, loan.ID,
	)
	return err
}

func (r *tidbRepository) SaveReward(ctx context.Context, reward *domain.Reward) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO rewards (id, user_id, asset_id, amount, reason, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		reward.ID, reward.UserID, reward.AssetID, reward.Amount, reward.Reason, reward.CreatedAt,
	)
	return err
}

// --- AppDataRepository Implementation ---

func (r *tidbRepository) GetAppData(ctx context.Context, id string) (*domain.AppData, error) {
	var appData domain.AppData
	err := r.db.QueryRowContext(ctx, "SELECT id, data FROM app_data WHERE id = ?", id).Scan(&appData.ID, &appData.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &appData, nil
}

func (r *tidbRepository) CreateAppData(ctx context.Context, id string, data []byte) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO app_data (id, data) VALUES (?, ?)", id, data)
	return err
}
