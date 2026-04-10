package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type GeometryRepository struct {
	db *sql.DB
}

func NewGeometryRepository(db *sql.DB) *GeometryRepository {
	return &GeometryRepository{db: db}
}

func (r *GeometryRepository) Create(ctx context.Context, geometry *domain.RoomGeometry) error {
	query := `
		INSERT INTO room_geometries (
			id, project_id, geometry_type, width, length, height, volume_m3,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		geometry.ID, geometry.ProjectID, geometry.GeometryType,
		geometry.Width, geometry.Length, geometry.Height, geometry.VolumeM3,
		geometry.CreatedAt, geometry.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create geometry: %w", err)
	}
	
	return nil
}

func (r *GeometryRepository) GetByProjectID(ctx context.Context, projectID string) (*domain.RoomGeometry, error) {
	query := `
		SELECT id, project_id, geometry_type, width, length, height, volume_m3,
			   created_at, updated_at
		FROM room_geometries 
		WHERE project_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	var geometry domain.RoomGeometry
	var geometryType string
	
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&geometry.ID, &geometry.ProjectID, &geometryType,
		&geometry.Width, &geometry.Length, &geometry.Height, &geometry.VolumeM3,
		&geometry.CreatedAt, &geometry.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get geometry: %w", err)
	}
	
	geometry.GeometryType = domain.GeometryType(geometryType)
	
	return &geometry, nil
}

func (r *GeometryRepository) GetByID(ctx context.Context, id string) (*domain.RoomGeometry, error) {
	query := `
		SELECT id, project_id, geometry_type, width, length, height, volume_m3,
			   created_at, updated_at
		FROM room_geometries 
		WHERE id = ?
	`
	
	var geometry domain.RoomGeometry
	var geometryType string
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&geometry.ID, &geometry.ProjectID, &geometryType,
		&geometry.Width, &geometry.Length, &geometry.Height, &geometry.VolumeM3,
		&geometry.CreatedAt, &geometry.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("geometry not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get geometry: %w", err)
	}
	
	geometry.GeometryType = domain.GeometryType(geometryType)
	
	return &geometry, nil
}

func (r *GeometryRepository) Update(ctx context.Context, geometry *domain.RoomGeometry) error {
	query := `
		UPDATE room_geometries 
		SET geometry_type = ?, width = ?, length = ?, height = ?, volume_m3 = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		geometry.GeometryType, geometry.Width, geometry.Length, geometry.Height,
		geometry.VolumeM3, geometry.UpdatedAt, geometry.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update geometry: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("geometry not found: %s", geometry.ID)
	}
	
	return nil
}

func (r *GeometryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM room_geometries WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete geometry: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("geometry not found: %s", id)
	}
	
	return nil
}

func (r *GeometryRepository) DeleteByProjectID(ctx context.Context, projectID string) error {
	query := `DELETE FROM room_geometries WHERE project_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete geometries by project: %w", err)
	}
	
	return nil
}
