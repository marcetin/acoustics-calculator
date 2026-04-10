package repo

import (
	"context"
	"database/sql"

	"acoustics-calculator/internal/domain"
)

type PlacementRepository struct {
	db *sql.DB
}

func NewPlacementRepository(db *sql.DB) *PlacementRepository {
	return &PlacementRepository{db: db}
}

func (r *PlacementRepository) Create(ctx context.Context, candidate *domain.PlacementCandidate) error {
	query := `
		INSERT INTO placement_candidates (id, project_id, analysis_run_id, surface_id, bounding_box_json,
			center_x, center_y, center_z, orientation, score, confidence, decision, diffuser_type_id,
			reasons_json, warnings_json, target_bands_json, source_receiver_pairs_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		candidate.ID,
		candidate.ProjectID,
		candidate.AnalysisRunID,
		candidate.SurfaceID,
		candidate.BoundingBoxJSON,
		candidate.CenterX,
		candidate.CenterY,
		candidate.CenterZ,
		candidate.Orientation,
		candidate.Score,
		candidate.Confidence,
		candidate.Decision,
		candidate.DiffuserTypeID,
		candidate.ReasonsJSON,
		candidate.WarningsJSON,
		candidate.TargetBandsJSON,
		candidate.SourceReceiverPairsJSON,
		candidate.CreatedAt,
		candidate.UpdatedAt,
	)
	return err
}

func (r *PlacementRepository) GetByID(ctx context.Context, id string) (*domain.PlacementCandidate, error) {
	query := `
		SELECT id, project_id, analysis_run_id, surface_id, bounding_box_json,
			center_x, center_y, center_z, orientation, score, confidence, decision, diffuser_type_id,
			reasons_json, warnings_json, target_bands_json, source_receiver_pairs_json, created_at, updated_at
		FROM placement_candidates
		WHERE id = ?
	`
	var c domain.PlacementCandidate
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.ProjectID,
		&c.AnalysisRunID,
		&c.SurfaceID,
		&c.BoundingBoxJSON,
		&c.CenterX,
		&c.CenterY,
		&c.CenterZ,
		&c.Orientation,
		&c.Score,
		&c.Confidence,
		&c.Decision,
		&c.DiffuserTypeID,
		&c.ReasonsJSON,
		&c.WarningsJSON,
		&c.TargetBandsJSON,
		&c.SourceReceiverPairsJSON,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PlacementRepository) ListByAnalysisRun(ctx context.Context, analysisRunID string) ([]*domain.PlacementCandidate, error) {
	query := `
		SELECT id, project_id, analysis_run_id, surface_id, bounding_box_json,
			center_x, center_y, center_z, orientation, score, confidence, decision, diffuser_type_id,
			reasons_json, warnings_json, target_bands_json, source_receiver_pairs_json, created_at, updated_at
		FROM placement_candidates
		WHERE analysis_run_id = ?
		ORDER BY score DESC
	`
	rows, err := r.db.QueryContext(ctx, query, analysisRunID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []*domain.PlacementCandidate
	for rows.Next() {
		var c domain.PlacementCandidate
		err := rows.Scan(
			&c.ID,
			&c.ProjectID,
			&c.AnalysisRunID,
			&c.SurfaceID,
			&c.BoundingBoxJSON,
			&c.CenterX,
			&c.CenterY,
			&c.CenterZ,
			&c.Orientation,
			&c.Score,
			&c.Confidence,
			&c.Decision,
			&c.DiffuserTypeID,
			&c.ReasonsJSON,
			&c.WarningsJSON,
			&c.TargetBandsJSON,
			&c.SourceReceiverPairsJSON,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, &c)
	}

	return candidates, nil
}

func (r *PlacementRepository) DeleteByAnalysisRun(ctx context.Context, analysisRunID string) error {
	query := `DELETE FROM placement_candidates WHERE analysis_run_id = ?`
	_, err := r.db.ExecContext(ctx, query, analysisRunID)
	return err
}

func (r *PlacementRepository) DeleteByProject(ctx context.Context, projectID string) error {
	query := `DELETE FROM placement_candidates WHERE project_id = ?`
	_, err := r.db.ExecContext(ctx, query, projectID)
	return err
}
