package repo

import (
	"context"
	"database/sql"

	"acoustics-calculator/internal/domain"
)

type DiffuserRepository struct {
	db *sql.DB
}

func NewDiffuserRepository(db *sql.DB) *DiffuserRepository {
	return &DiffuserRepository{db: db}
}

func (r *DiffuserRepository) ListPresets(ctx context.Context) ([]domain.DiffuserType, error) {
	query := `
		SELECT id, name, category, panel_width_m, panel_height_m, panel_depth_m,
			   min_effective_hz, max_effective_hz, scatter_axis, mount_orientation_json,
			   cost_class, is_preset, created_at, updated_at
		FROM diffuser_types
		WHERE is_preset = 1
		ORDER BY name
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var diffusers []domain.DiffuserType
	for rows.Next() {
		var d domain.DiffuserType
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Category,
			&d.PanelWidthM,
			&d.PanelHeightM,
			&d.PanelDepthM,
			&d.MinEffectiveHz,
			&d.MaxEffectiveHz,
			&d.ScatterAxis,
			&d.MountOrientationJSON,
			&d.CostClass,
			&d.IsPreset,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		diffusers = append(diffusers, d)
	}

	return diffusers, nil
}

func (r *DiffuserRepository) GetByID(ctx context.Context, id string) (*domain.DiffuserType, error) {
	query := `
		SELECT id, name, category, panel_width_m, panel_height_m, panel_depth_m,
			   min_effective_hz, max_effective_hz, scatter_axis, mount_orientation_json,
			   cost_class, is_preset, created_at, updated_at
		FROM diffuser_types
		WHERE id = ?
	`
	var d domain.DiffuserType
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID,
		&d.Name,
		&d.Category,
		&d.PanelWidthM,
		&d.PanelHeightM,
		&d.PanelDepthM,
		&d.MinEffectiveHz,
		&d.MaxEffectiveHz,
		&d.ScatterAxis,
		&d.MountOrientationJSON,
		&d.CostClass,
		&d.IsPreset,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DiffuserRepository) Create(ctx context.Context, diffuser *domain.DiffuserType) error {
	query := `
		INSERT INTO diffuser_types (id, name, category, panel_width_m, panel_height_m, panel_depth_m,
			min_effective_hz, max_effective_hz, scatter_axis, mount_orientation_json, cost_class, is_preset)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		diffuser.ID,
		diffuser.Name,
		diffuser.Category,
		diffuser.PanelWidthM,
		diffuser.PanelHeightM,
		diffuser.PanelDepthM,
		diffuser.MinEffectiveHz,
		diffuser.MaxEffectiveHz,
		diffuser.ScatterAxis,
		diffuser.MountOrientationJSON,
		diffuser.CostClass,
		diffuser.IsPreset,
	)
	return err
}
