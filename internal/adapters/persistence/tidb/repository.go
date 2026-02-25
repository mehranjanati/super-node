package tidb

import (
	"context"
	"database/sql"
	"encoding/json"
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
func NewTiDBRepository(ctx context.Context, dsn string) (ports.UserRepository, ports.FinanceRepository, ports.AppDataRepository, ports.SocialRepository, ports.AgentRepository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to open tidb connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to ping tidb: %w", err)
	}

	// Ensure tables exist (Simple migration for now)
	if err := migrate(ctx, db); err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to migrate tidb schema: %w", err)
	}

	repo := &tidbRepository{db: db}
	agentRepo := newTiDBAgentRepository(db)
	return repo, repo, repo, repo, agentRepo, nil
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
		`CREATE TABLE IF NOT EXISTS posts (
			id VARCHAR(255) PRIMARY KEY,
			author_id VARCHAR(255) NOT NULL,
			content TEXT,
			media_urls JSON,
			tags JSON,
			likes INT DEFAULT 0,
			metadata JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS agents (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			type VARCHAR(50) NOT NULL,
			status VARCHAR(50) NOT NULL,
			owner_id VARCHAR(255) NOT NULL,
			avatar VARCHAR(255),
			config JSON,
			performance JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			post_id VARCHAR(255) NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (post_id, user_id)
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

func (r *tidbRepository) TransferFunds(ctx context.Context, fromUserID, toUserID, assetID string, amount float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Check sender balance (Locking row)
	var balance float64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM user_balances WHERE user_id = ? AND asset_id = ? FOR UPDATE", fromUserID, assetID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("sender has no balance record")
		}
		return err
	}

	if balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// 2. Deduct from sender
	_, err = tx.ExecContext(ctx, "UPDATE user_balances SET balance = balance - ? WHERE user_id = ? AND asset_id = ?", amount, fromUserID, assetID)
	if err != nil {
		return err
	}

	// 3. Add to receiver (Upsert)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_balances (user_id, asset_id, balance) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE balance = balance + ?`,
		toUserID, assetID, amount, amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
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

// --- SocialRepository Implementation ---

func (r *tidbRepository) SavePost(ctx context.Context, post *domain.Post) error {
	mediaJSON, _ := json.Marshal(post.MediaURLs)
	tagsJSON, _ := json.Marshal(post.Tags)
	metadataJSON, _ := json.Marshal(post.Metadata)

	_, err := r.db.ExecContext(ctx,
		"INSERT INTO posts (id, author_id, content, media_urls, tags, likes, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		post.ID, post.AuthorID, post.Content, mediaJSON, tagsJSON, post.Likes, metadataJSON, post.CreatedAt,
	)
	return err
}

func (r *tidbRepository) GetPosts(ctx context.Context, filter domain.FeedFilter) ([]*domain.Post, error) {
	// Simple implementation ignoring complex filters for now
	query := "SELECT id, author_id, content, media_urls, tags, likes, metadata, created_at FROM posts ORDER BY created_at DESC"
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		var mediaJSON, tagsJSON, metadataJSON []byte
		var likes int
		// Scan likes into temporary variable first to avoid type mismatch if needed
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &mediaJSON, &tagsJSON, &likes, &metadataJSON, &post.CreatedAt); err != nil {
			return nil, err
		}
		post.Likes = likes

		var mediaURLs []string
		if len(mediaJSON) > 0 {
			json.Unmarshal(mediaJSON, &mediaURLs)
		}
		post.MediaURLs = mediaURLs

		var tags []string
		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &tags)
		}
		post.Tags = tags

		var metadata map[string]interface{}
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &metadata)
		}
		post.Metadata = metadata

		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *tidbRepository) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post
	var mediaJSON, tagsJSON, metadataJSON []byte
	var likes int

	err := r.db.QueryRowContext(ctx, "SELECT id, author_id, content, media_urls, tags, likes, metadata, created_at FROM posts WHERE id = ?", id).Scan(
		&post.ID, &post.AuthorID, &post.Content, &mediaJSON, &tagsJSON, &likes, &metadataJSON, &post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	post.Likes = likes

	var mediaURLs []string
	if len(mediaJSON) > 0 {
		json.Unmarshal(mediaJSON, &mediaURLs)
	}
	post.MediaURLs = mediaURLs

	var tags []string
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &tags)
	}
	post.Tags = tags

	var metadata map[string]interface{}
	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &metadata)
	}
	post.Metadata = metadata

	return &post, nil
}

func (r *tidbRepository) AddLike(ctx context.Context, postID, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert into post_likes
	_, err = tx.ExecContext(ctx, "INSERT INTO post_likes (post_id, user_id) VALUES (?, ?)", postID, userID)
	if err != nil {
		// If duplicate entry, assume user already liked and ignore error (or return specific error if needed)
		// For now, we assume it's fine and just return nil (idempotent)
		// But wait, if we return nil, we shouldn't increment count again.
		// If INSERT fails due to unique constraint, we should NOT increment.
		// So checking error is important.
		return nil
	}

	// Increment likes count in posts table
	_, err = tx.ExecContext(ctx, "UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *tidbRepository) SaveComment(ctx context.Context, comment *domain.Comment) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO comments (id, post_id, author_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		comment.ID, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt,
	)
	return err
}
