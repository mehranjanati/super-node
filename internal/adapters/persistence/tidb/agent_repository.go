package tidb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"nexus-super-node-v3/internal/core/domain"
	"nexus-super-node-v3/internal/ports"
	"time"
)

type tidbAgentRepository struct {
	db *sql.DB
}

// NewTiDBAgentRepository is internal now, used by main repo constructor
func newTiDBAgentRepository(db *sql.DB) ports.AgentRepository {
	return &tidbAgentRepository{
		db: db,
	}
}

func (r *tidbAgentRepository) Create(ctx context.Context, agent *domain.Agent) error {
	configJSON, _ := json.Marshal(agent.Config)
	perfJSON, _ := json.Marshal(agent.Performance)

	query := `
		INSERT INTO agents (id, name, description, type, status, owner_id, avatar, config, performance, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		agent.ID, agent.Name, agent.Description, agent.Type, agent.Status, agent.OwnerID, agent.Avatar,
		string(configJSON), string(perfJSON), agent.CreatedAt, agent.UpdatedAt)
	return err
}

func (r *tidbAgentRepository) GetByID(ctx context.Context, id string) (*domain.Agent, error) {
	query := `SELECT id, name, description, type, status, owner_id, avatar, config, performance, created_at, updated_at FROM agents WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var agent domain.Agent
	var configJSON, perfJSON []byte

	// TiDB/MySQL timestamps might come as []uint8 if parseTime=true is not set,
	// but assuming standard usage. If issues arise, we can scan to string.
	err := row.Scan(
		&agent.ID, &agent.Name, &agent.Description, &agent.Type, &agent.Status, &agent.OwnerID, &agent.Avatar,
		&configJSON, &perfJSON, &agent.CreatedAt, &agent.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found")
		}
		return nil, err
	}

	if len(configJSON) > 0 {
		json.Unmarshal(configJSON, &agent.Config)
	}
	if len(perfJSON) > 0 {
		json.Unmarshal(perfJSON, &agent.Performance)
	}

	return &agent, nil
}

func (r *tidbAgentRepository) List(ctx context.Context, ownerID string) ([]*domain.Agent, error) {
	query := `SELECT id, name, description, type, status, owner_id, avatar, config, performance, created_at, updated_at FROM agents`
	args := []interface{}{}

	if ownerID != "" {
		query += ` WHERE owner_id = ?`
		args = append(args, ownerID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agents := []*domain.Agent{}
	for rows.Next() {
		var agent domain.Agent
		var configJSON, perfJSON []byte

		err := rows.Scan(
			&agent.ID, &agent.Name, &agent.Description, &agent.Type, &agent.Status, &agent.OwnerID, &agent.Avatar,
			&configJSON, &perfJSON, &agent.CreatedAt, &agent.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if len(configJSON) > 0 {
			json.Unmarshal(configJSON, &agent.Config)
		}
		if len(perfJSON) > 0 {
			json.Unmarshal(perfJSON, &agent.Performance)
		}
		agents = append(agents, &agent)
	}
	return agents, nil
}

func (r *tidbAgentRepository) Update(ctx context.Context, agent *domain.Agent) error {
	configJSON, _ := json.Marshal(agent.Config)
	perfJSON, _ := json.Marshal(agent.Performance)

	query := `
		UPDATE agents 
		SET name=?, description=?, type=?, status=?, avatar=?, config=?, performance=?, updated_at=?
		WHERE id=?
	`
	_, err := r.db.ExecContext(ctx, query,
		agent.Name, agent.Description, agent.Type, agent.Status, agent.Avatar,
		string(configJSON), string(perfJSON), time.Now(), agent.ID)
	return err
}

func (r *tidbAgentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM agents WHERE id = ?", id)
	return err
}

func (r *tidbAgentRepository) UpdateStatus(ctx context.Context, id string, status domain.AgentStatus) error {
	_, err := r.db.ExecContext(ctx, "UPDATE agents SET status = ?, updated_at = ? WHERE id = ?", status, time.Now(), id)
	return err
}
