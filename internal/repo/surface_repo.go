package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type SurfaceRepository struct {
	db *sql.DB
}

func NewSurfaceRepository(db *sql.DB) *SurfaceRepository {
	return &SurfaceRepository{db: db}
}

func (r *SurfaceRepository) Create(ctx context.Context, surface *domain.Surface) error {
	query := `
		INSERT INTO surfaces (
			id, project_id, geometry_id, name, kind, polygon_3d_json,
			normal_x, normal_y, normal_z, area_m2, material_id, is_mountable,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		surface.ID, surface.ProjectID, surface.GeometryID, surface.Name, surface.Kind,
		surface.Polygon3DJSON, surface.NormalX, surface.NormalY, surface.NormalZ,
		surface.AreaM2, surface.MaterialID, surface.IsMountable,
		surface.CreatedAt, surface.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create surface: %w", err)
	}
	
	return nil
}

func (r *SurfaceRepository) GetByID(ctx context.Context, id string) (*domain.Surface, error) {
	query := `
		SELECT id, project_id, geometry_id, name, kind, polygon_3d_json,
			   normal_x, normal_y, normal_z, area_m2, material_id, is_mountable,
			   created_at, updated_at
		FROM surfaces 
		WHERE id = ?
	`
	
	var surface domain.Surface
	var kind string
	var materialID sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&surface.ID, &surface.ProjectID, &surface.GeometryID, &surface.Name, &kind,
		&surface.Polygon3DJSON, &surface.NormalX, &surface.NormalY, &surface.NormalZ,
		&surface.AreaM2, &materialID, &surface.IsMountable,
		&surface.CreatedAt, &surface.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("surface not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get surface: %w", err)
	}
	
	surface.Kind = domain.SurfaceKind(kind)
	if materialID.Valid {
		surface.MaterialID = &materialID.String
	}
	
	return &surface, nil
}

func (r *SurfaceRepository) ListByProjectID(ctx context.Context, projectID string) ([]*domain.Surface, error) {
	query := `
		SELECT id, project_id, geometry_id, name, kind, polygon_3d_json,
			   normal_x, normal_y, normal_z, area_m2, material_id, is_mountable,
			   created_at, updated_at
		FROM surfaces 
		WHERE project_id = ?
		ORDER BY kind, name
	`
	
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list surfaces: %w", err)
	}
	defer rows.Close()
	
	var surfaces []*domain.Surface
	
	for rows.Next() {
		var surface domain.Surface
		var kind string
		var materialID sql.NullString
		
		err := rows.Scan(
			&surface.ID, &surface.ProjectID, &surface.GeometryID, &surface.Name, &kind,
			&surface.Polygon3DJSON, &surface.NormalX, &surface.NormalY, &surface.NormalZ,
			&surface.AreaM2, &materialID, &surface.IsMountable,
			&surface.CreatedAt, &surface.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan surface: %w", err)
		}
		
		surface.Kind = domain.SurfaceKind(kind)
		if materialID.Valid {
			surface.MaterialID = &materialID.String
		}
		
		surfaces = append(surfaces, &surface)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating surfaces: %w", err)
	}
	
	return surfaces, nil
}

func (r *SurfaceRepository) Update(ctx context.Context, surface *domain.Surface) error {
	query := `
		UPDATE surfaces 
		SET name = ?, material_id = ?, is_mountable = ?, updated_at = ?
		WHERE id = ?
	`
	
	var materialID interface{} = nil
	if surface.MaterialID != nil {
		materialID = *surface.MaterialID
	}
	
	result, err := r.db.ExecContext(ctx, query,
		surface.Name, materialID, surface.IsMountable,
		surface.UpdatedAt, surface.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update surface: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("surface not found: %s", surface.ID)
	}
	
	return nil
}

func (r *SurfaceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM surfaces WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete surface: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("surface not found: %s", id)
	}
	
	return nil
}

func (r *SurfaceRepository) DeleteByProjectID(ctx context.Context, projectID string) error {
	query := `DELETE FROM surfaces WHERE project_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete surfaces by project: %w", err)
	}
	
	return nil
}

func (r *SurfaceRepository) DeleteByGeometryID(ctx context.Context, geometryID string) error {
	query := `DELETE FROM surfaces WHERE geometry_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, geometryID)
	if err != nil {
		return fmt.Errorf("failed to delete surfaces by geometry: %w", err)
	}
	
	return nil
}
