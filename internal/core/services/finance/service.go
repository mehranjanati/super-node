package finance

import (
	"context"
	"fmt"
	"time"

	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"

	"github.com/google/uuid"
)

type FinanceServiceImpl struct {
	repo ports.FinanceRepository
}

func NewFinanceService(repo ports.FinanceRepository) ports.FinanceService {
	return &FinanceServiceImpl{repo: repo}
}

func (s *FinanceServiceImpl) CreateAsset(ctx context.Context, name, symbol string, assetType domain.AssetType, initialSupply float64) (*domain.Asset, error) {
	asset := &domain.Asset{
		ID:          uuid.New().String(),
		Name:        name,
		Symbol:      symbol,
		Type:        assetType,
		TotalSupply: initialSupply,
		OwnerID:     "system", // Or pass ownerID
		CreatedAt:   time.Now(),
	}

	if err := s.repo.SaveAsset(ctx, asset); err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return asset, nil
}

func (s *FinanceServiceImpl) GetAsset(ctx context.Context, id string) (*domain.Asset, error) {
	return s.repo.GetAsset(ctx, id)
}

func (s *FinanceServiceImpl) GetUserBalance(ctx context.Context, userID, assetID string) (*domain.UserBalance, error) {
	return s.repo.GetBalance(ctx, userID, assetID)
}

func (s *FinanceServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID, assetID string, amount float64) error {
	// Simple transfer logic (without transaction for brevity, but should be transactional)
	senderBalance, err := s.repo.GetBalance(ctx, fromUserID, assetID)
	if err != nil {
		return err
	}

	if senderBalance.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	if err := s.repo.UpdateUserBalance(ctx, fromUserID, assetID, -amount); err != nil {
		return err
	}

	if err := s.repo.UpdateUserBalance(ctx, toUserID, assetID, amount); err != nil {
		// Rollback sender would be needed here in a real TX
		return err
	}

	return nil
}

func (s *FinanceServiceImpl) RequestLoan(ctx context.Context, borrowerID, collateralID string, amount float64) (*domain.Loan, error) {
	// Check collateral logic...
	loan := &domain.Loan{
		ID:           uuid.New().String(),
		BorrowerID:   borrowerID,
		CollateralID: collateralID,
		Amount:       amount,
		InterestRate: 0.05, // 5% fixed for now
		Status:       "ACTIVE",
		DueDate:      time.Now().AddDate(0, 1, 0), // 1 month
	}

	if err := s.repo.SaveLoan(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *FinanceServiceImpl) RepayLoan(ctx context.Context, loanID string) error {
	// Fetch loan, check balance, etc.
	loan := &domain.Loan{ID: loanID, Status: "REPAID"}
	return s.repo.UpdateLoan(ctx, loan)
}

func (s *FinanceServiceImpl) DistributeReward(ctx context.Context, userID, assetID string, amount float64, reason string) error {
	reward := &domain.Reward{
		ID:        uuid.New().String(),
		UserID:    userID,
		AssetID:   assetID,
		Amount:    amount,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	if err := s.repo.SaveReward(ctx, reward); err != nil {
		return err
	}

	return s.repo.UpdateUserBalance(ctx, userID, assetID, amount)
}
