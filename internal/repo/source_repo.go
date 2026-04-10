package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type SourceRepository struct {
	db *sql.DB
}

func NewSourceRepository(db *sql.DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func (r *SourceRepository) Create(ctx context.Context, source *domain.Source) error {
	query := `
		INSERT INTO sources (
			id, project_id, name, type, position_x, position_y, position_z,
			aim_x, aim_y, aim_z, gain_db, group_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	var groupID interface{} = nil
	if source.GroupID != nil {
		groupID = *source.GroupID
	}
	
	_, err := r.db.ExecContext(ctx, query,
		source.ID, source.ProjectID, source.Name, source.Type,
		source.PositionX, source.PositionY, source.PositionZ,
		source.AimX, source.AimY, source.AimZ, source.GainDB,
		groupID, source.CreatedAt, source.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}
	
	return nil
}

func (r *SourceRepository) GetByID(ctx context.Context, id string) (*domain.Source, error) {
	query := `
		SELECT id, project_id, name, type, position_x, position_y, position_z,
			   aim_x, aim_y, aim_z, gain_db, group_id, created_at, updated_at
		FROM sources 
		WHERE id = ?
	`
	
	var source domain.Source
	var sourceType string
	var groupID sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&source.ID, &source.ProjectID, &source.Name, &sourceType,
		&source.PositionX, &source.PositionY, &source.PositionZ,
		&source.AimX, &source.AimY, &source.AimZ, &source.GainDB,
		&groupID, &source.CreatedAt, &source.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("source not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get source: %w", err)
	}
	
	source.Type = domain.SourceType(sourceType)
	if groupID.Valid {
		source.GroupID = &groupID.String
	}
	
	return &source, nil
}

func (r *SourceRepository) ListByProjectID(ctx context.Context, projectID string) ([]*domain.Source, error) {
	query := `
		SELECT id, project_id, name, type, position_x, position_y, position_z,
			   aim_x, aim_y, aim_z, gain_db, group_id, created_at, updated_at
		FROM sources 
		WHERE project_id = ?
		ORDER BY created_at
	`
	
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}
	defer rows.Close()
	
	var sources []*domain.Source
	
	for rows.Next() {
		var source domain.Source
		var sourceType string
		var groupID sql.NullString
		
		err := rows.Scan(
			&source.ID, &source.ProjectID, &source.Name, &sourceType,
			&source.PositionX, &source.PositionY, &source.PositionZ,
			&source.AimX, &source.AimY, &source.AimZ, &source.GainDB,
			&groupID, &source.CreatedAt, &source.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		
		source.Type = domain.SourceType(sourceType)
		if groupID.Valid {
			source.GroupID = &groupID.String
		}
		
		sources = append(sources, &source)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sources: %w", err)
	}
	
	return sources, nil
}

func (r *SourceRepository) Update(ctx context.Context, source *domain.Source) error {
	query := `
		UPDATE sources 
		SET name = ?, type = ?, position_x = ?, position_y = ?, position_z = ?,
		    aim_x = ?, aim_y = ?, aim_z = ?, gain_db = ?, group_id = ?, updated_at = ?
		WHERE id = ?
	`
	
	var groupID interface{} = nil
	if source.GroupID != nil {
		groupID = *source.GroupID
	}
	
	result, err := r.db.ExecContext(ctx, query,
		source.Name, source.Type, source.PositionX, source.PositionY, source.PositionZ,
		source.AimX, source.AimY, source.AimZ, source.GainDB, groupID,
		source.UpdatedAt, source.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update source: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("source not found: %s", source.ID)
	}
	
	return nil
}

func (r *SourceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sources WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("source not found: %s", id)
	}
	
	return nil
}

func (r *SourceRepository) DeleteByProjectID(ctx context.Context, projectID string) error {
	query := `DELETE FROM sources WHERE project_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete sources by project: %w", err)
	}
	
	return nil
}
