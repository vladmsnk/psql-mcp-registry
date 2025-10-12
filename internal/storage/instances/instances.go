package instances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	"psql-mcp-registry/internal/model"
)

var ErrNotFound = errors.New("instance not found")

func (s *PostgresStorage) CreateInstance(ctx context.Context, instance *model.Instance) error {
	query := `
		INSERT INTO instance_registry 
		(name, database_name, description, creator_username, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx, query,
		instance.Name,
		instance.DatabaseName,
		instance.Description,
		instance.CreatorUsername,
		instance.Status,
	).Scan(&instance.ID, &instance.CreatedAt, &instance.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetInstanceByName(ctx context.Context, name string) (*model.Instance, error) {
	query := `
		SELECT 
			id, name, database_name, description, creator_username, 
			status, created_at, updated_at
		FROM instance_registry
		WHERE name = $1
	`

	var instance model.Instance
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&instance.ID,
		&instance.Name,
		&instance.DatabaseName,
		&instance.Description,
		&instance.CreatorUsername,
		&instance.Status,
		&instance.CreatedAt,
		&instance.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return &instance, nil
}

func (s *PostgresStorage) ListInstances(ctx context.Context) ([]model.Instance, error) {
	query := `
		SELECT 
			id, name, database_name, description, creator_username,
			status, created_at, updated_at
		FROM instance_registry
		WHERE status = 'active'
		ORDER BY name
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}
	defer rows.Close()

	var instances []model.Instance
	for rows.Next() {
		var inst model.Instance
		err := rows.Scan(
			&inst.ID,
			&inst.Name,
			&inst.DatabaseName,
			&inst.Description,
			&inst.CreatorUsername,
			&inst.Status,
			&inst.CreatedAt,
			&inst.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan instance: %w", err)
		}
		instances = append(instances, inst)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return instances, nil
}
