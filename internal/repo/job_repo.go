package repo

import (
	"context"
	"database/sql"

	"acoustics-calculator/internal/domain"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	query := `
		INSERT INTO jobs (id, project_id, kind, status, phase, progress_percent, message, error_message, created_at, started_at, finished_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		job.ID,
		job.ProjectID,
		job.Kind,
		job.Status,
		job.Phase,
		job.ProgressPercent,
		job.Message,
		job.ErrorMessage,
		job.CreatedAt,
		job.StartedAt,
		job.FinishedAt,
	)
	return err
}

func (r *JobRepository) Update(ctx context.Context, job *domain.Job) error {
	query := `
		UPDATE jobs
		SET status = ?, phase = ?, progress_percent = ?, message = ?, error_message = ?, started_at = ?, finished_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		job.Status,
		job.Phase,
		job.ProgressPercent,
		job.Message,
		job.ErrorMessage,
		job.StartedAt,
		job.FinishedAt,
		job.ID,
	)
	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (*domain.Job, error) {
	query := `
		SELECT id, project_id, kind, status, phase, progress_percent, message, error_message, created_at, started_at, finished_at
		FROM jobs
		WHERE id = ?
	`
	var job domain.Job
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.ProjectID,
		&job.Kind,
		&job.Status,
		&job.Phase,
		&job.ProgressPercent,
		&job.Message,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.StartedAt,
		&job.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) GetLatestByProject(ctx context.Context, projectID string) (*domain.Job, error) {
	query := `
		SELECT id, project_id, kind, status, phase, progress_percent, message, error_message, created_at, started_at, finished_at
		FROM jobs
		WHERE project_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var job domain.Job
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&job.ID,
		&job.ProjectID,
		&job.Kind,
		&job.Status,
		&job.Phase,
		&job.ProgressPercent,
		&job.Message,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.StartedAt,
		&job.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) GetLatestActiveByProject(ctx context.Context, projectID string) (*domain.Job, error) {
	query := `
		SELECT id, project_id, kind, status, phase, progress_percent, message, error_message, created_at, started_at, finished_at
		FROM jobs
		WHERE project_id = ? AND status IN ('QUEUED', 'RUNNING')
		ORDER BY created_at DESC
		LIMIT 1
	`
	var job domain.Job
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&job.ID,
		&job.ProjectID,
		&job.Kind,
		&job.Status,
		&job.Phase,
		&job.ProgressPercent,
		&job.Message,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.StartedAt,
		&job.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) GetLatestByProjectAndKind(ctx context.Context, projectID string, kind domain.JobKind) (*domain.Job, error) {
	query := `
		SELECT id, project_id, kind, status, phase, progress_percent, message, error_message, created_at, started_at, finished_at
		FROM jobs
		WHERE project_id = ? AND kind = ?
		ORDER BY created_at DESC
		LIMIT 1
	`
	var job domain.Job
	err := r.db.QueryRowContext(ctx, query, projectID, kind).Scan(
		&job.ID,
		&job.ProjectID,
		&job.Kind,
		&job.Status,
		&job.Phase,
		&job.ProgressPercent,
		&job.Message,
		&job.ErrorMessage,
		&job.CreatedAt,
		&job.StartedAt,
		&job.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) DeleteByProject(ctx context.Context, projectID string) error {
	query := `DELETE FROM jobs WHERE project_id = ?`
	_, err := r.db.ExecContext(ctx, query, projectID)
	return err
}
