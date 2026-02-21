package ports

import (
	"context"
	"nexus-super-node-v3/internal/core/domain"
)

// FinanceService defines the business logic for finance and tokenization
type FinanceService interface {
	// Asset Management
	CreateAsset(ctx context.Context, name, symbol string, assetType domain.AssetType, initialSupply float64) (*domain.Asset, error)
	GetAsset(ctx context.Context, id string) (*domain.Asset, error)
	
	// Balance Management
	GetUserBalance(ctx context.Context, userID, assetID string) (*domain.UserBalance, error)
	Transfer(ctx context.Context, fromUserID, toUserID, assetID string, amount float64) error

	// DeFi
	RequestLoan(ctx context.Context, borrowerID, collateralID string, amount float64) (*domain.Loan, error)
	RepayLoan(ctx context.Context, loanID string) error
	
	// Rewards
	DistributeReward(ctx context.Context, userID, assetID string, amount float64, reason string) error
}

// FinanceRepository defines the persistence interface for finance data
type FinanceRepository interface {
	SaveAsset(ctx context.Context, asset *domain.Asset) error
	GetAsset(ctx context.Context, id string) (*domain.Asset, error)
	UpdateUserBalance(ctx context.Context, userID, assetID string, delta float64) error
	GetBalance(ctx context.Context, userID, assetID string) (*domain.UserBalance, error)
	
	SaveLoan(ctx context.Context, loan *domain.Loan) error
	UpdateLoan(ctx context.Context, loan *domain.Loan) error
	
	SaveReward(ctx context.Context, reward *domain.Reward) error
}
