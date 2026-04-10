package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type MaterialRepository struct {
	db *sql.DB
}

func NewMaterialRepository(db *sql.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) Create(ctx context.Context, material *domain.Material) error {
	query := `
		INSERT INTO materials (
			id, project_id, name, category, absorption_bands_json,
			scattering_bands_json, is_preset, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	var projectID interface{} = nil
	if material.ProjectID != nil {
		projectID = *material.ProjectID
	}
	
	_, err := r.db.ExecContext(ctx, query,
		material.ID, projectID, material.Name, material.Category,
		material.AbsorptionBandsJSON, material.ScatteringBandsJSON,
		material.IsPreset, material.CreatedAt, material.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create material: %w", err)
	}
	
	return nil
}

func (r *MaterialRepository) GetByID(ctx context.Context, id string) (*domain.Material, error) {
	query := `
		SELECT id, project_id, name, category, absorption_bands_json,
			   scattering_bands_json, is_preset, created_at, updated_at
		FROM materials 
		WHERE id = ?
	`
	
	var material domain.Material
	var projectID sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&material.ID, &projectID, &material.Name, &material.Category,
		&material.AbsorptionBandsJSON, &material.ScatteringBandsJSON,
		&material.IsPreset, &material.CreatedAt, &material.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("material not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get material: %w", err)
	}
	
	if projectID.Valid {
		material.ProjectID = &projectID.String
	}
	
	return &material, nil
}

func (r *MaterialRepository) ListPresets(ctx context.Context) ([]*domain.Material, error) {
	query := `
		SELECT id, project_id, name, category, absorption_bands_json,
			   scattering_bands_json, is_preset, created_at, updated_at
		FROM materials 
		WHERE is_preset = 1
		ORDER BY category, name
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list preset materials: %w", err)
	}
	defer rows.Close()
	
	var materials []*domain.Material
	
	for rows.Next() {
		var material domain.Material
		var projectID sql.NullString
		
		err := rows.Scan(
			&material.ID, &projectID, &material.Name, &material.Category,
			&material.AbsorptionBandsJSON, &material.ScatteringBandsJSON,
			&material.IsPreset, &material.CreatedAt, &material.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan material: %w", err)
		}
		
		if projectID.Valid {
			material.ProjectID = &projectID.String
		}
		
		materials = append(materials, &material)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating materials: %w", err)
	}
	
	return materials, nil
}

func (r *MaterialRepository) ListByProjectID(ctx context.Context, projectID string) ([]*domain.Material, error) {
	query := `
		SELECT id, project_id, name, category, absorption_bands_json,
			   scattering_bands_json, is_preset, created_at, updated_at
		FROM materials 
		WHERE project_id = ? OR project_id IS NULL
		ORDER BY is_preset DESC, category, name
	`
	
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list project materials: %w", err)
	}
	defer rows.Close()
	
	var materials []*domain.Material
	
	for rows.Next() {
		var material domain.Material
		var pid sql.NullString
		
		err := rows.Scan(
			&material.ID, &pid, &material.Name, &material.Category,
			&material.AbsorptionBandsJSON, &material.ScatteringBandsJSON,
			&material.IsPreset, &material.CreatedAt, &material.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan material: %w", err)
		}
		
		if pid.Valid {
			material.ProjectID = &pid.String
		}
		
		materials = append(materials, &material)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating materials: %w", err)
	}
	
	return materials, nil
}

func (r *MaterialRepository) Update(ctx context.Context, material *domain.Material) error {
	query := `
		UPDATE materials 
		SET name = ?, category = ?, absorption_bands_json = ?, 
		    scattering_bands_json = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		material.Name, material.Category, material.AbsorptionBandsJSON,
		material.ScatteringBandsJSON, material.UpdatedAt, material.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("material not found: %s", material.ID)
	}
	
	return nil
}

func (r *MaterialRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM materials WHERE id = ? AND is_preset = 0`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("material not found or is preset: %s", id)
	}
	
	return nil
}
