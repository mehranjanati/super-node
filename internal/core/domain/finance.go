package domain

import "time"

// AssetType represents the type of asset (e.g., Token, Equity, NFT)
type AssetType string

const (
	AssetTypeToken  AssetType = "TOKEN"
	AssetTypeEquity AssetType = "EQUITY"
	AssetTypeNFT    AssetType = "NFT"
)

// Asset represents a tokenized asset on the platform
type Asset struct {
	ID          string    `json:"id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Type        AssetType `json:"type"`
	TotalSupply float64   `json:"total_supply"`
	OwnerID     string    `json:"owner_id"` // Organization or User who issued it
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserBalance represents a user's holding of an asset
type UserBalance struct {
	UserID  string  `json:"user_id"`
	AssetID string  `json:"asset_id"`
	Balance float64 `json:"balance"`
}

// Loan represents a DeFi loan record
type Loan struct {
	ID           string    `json:"id"`
	BorrowerID   string    `json:"borrower_id"`
	CollateralID string    `json:"collateral_id"` // Asset ID used as collateral
	Amount       float64   `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	Status       string    `json:"status"` // ACTIVE, REPAID, LIQUIDATED
	DueDate      time.Time `json:"due_date"`
}

// Reward represents a distribution of tokens to users
type Reward struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	AssetID   string    `json:"asset_id"`
	Amount    float64   `json:"amount"`
	Reason    string    `json:"reason"` // e.g., "staking", "referral"
	CreatedAt time.Time `json:"created_at"`
}
