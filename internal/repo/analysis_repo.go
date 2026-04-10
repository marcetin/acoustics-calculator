package repo

import (
	"context"
	"database/sql"

	"acoustics-calculator/internal/domain"
)

type AnalysisRepository struct {
	db *sql.DB
}

func NewAnalysisRepository(db *sql.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Create(ctx context.Context, analysis *domain.AnalysisRun) error {
	query := `
		INSERT INTO analysis_runs (id, project_id, status, inputs_hash, metrics_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		analysis.ID,
		analysis.ProjectID,
		analysis.Status,
		analysis.InputsHash,
		analysis.MetricsJSON,
		analysis.CreatedAt,
	)
	return err
}

func (r *AnalysisRepository) Get(ctx context.Context, id string) (*domain.AnalysisRun, error) {
	query := `
		SELECT id, project_id, status, inputs_hash, metrics_json, started_at, finished_at, error_message, created_at
		FROM analysis_runs
		WHERE id = ?
	`
	var analysis domain.AnalysisRun
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&analysis.ID,
		&analysis.ProjectID,
		&analysis.Status,
		&analysis.InputsHash,
		&analysis.MetricsJSON,
		&analysis.StartedAt,
		&analysis.FinishedAt,
		&analysis.ErrorMessage,
		&analysis.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *AnalysisRepository) GetLatestByProject(ctx context.Context, projectID string) (*domain.AnalysisRun, error) {
	query := `
		SELECT id, project_id, status, inputs_hash, metrics_json, started_at, finished_at, error_message, created_at
		FROM analysis_runs
		WHERE project_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var analysis domain.AnalysisRun
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&analysis.ID,
		&analysis.ProjectID,
		&analysis.Status,
		&analysis.InputsHash,
		&analysis.MetricsJSON,
		&analysis.StartedAt,
		&analysis.FinishedAt,
		&analysis.ErrorMessage,
		&analysis.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *AnalysisRepository) Update(ctx context.Context, analysis *domain.AnalysisRun) error {
	query := `
		UPDATE analysis_runs
		SET status = ?, metrics_json = ?, started_at = ?, finished_at = ?, error_message = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		analysis.Status,
		analysis.MetricsJSON,
		analysis.StartedAt,
		analysis.FinishedAt,
		analysis.ErrorMessage,
		analysis.ID,
	)
	return err
}

func (r *AnalysisRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM analysis_runs WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
