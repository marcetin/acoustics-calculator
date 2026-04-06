package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	query := `
		INSERT INTO projects (
			id, name, space_type, goal_type, units, 
			temperature_c, humidity_percent, speed_of_sound,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		project.ID, project.Name, project.SpaceType, project.GoalType, project.Units,
		project.TemperatureC, project.HumidityPercent, project.SpeedOfSound,
		project.CreatedAt, project.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `
		SELECT id, name, space_type, goal_type, units,
			   temperature_c, humidity_percent, speed_of_sound,
			   created_at, updated_at
		FROM projects 
		WHERE id = ?
	`

	var project domain.Project
	var spaceType, goalType, units string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&project.ID, &project.Name, &spaceType, &goalType, &units,
		&project.TemperatureC, &project.HumidityPercent, &project.SpeedOfSound,
		&project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	project.SpaceType = domain.SpaceType(spaceType)
	project.GoalType = domain.GoalType(goalType)
	project.Units = domain.Units(units)

	return &project, nil
}

func (r *ProjectRepository) List(ctx context.Context) ([]*domain.Project, error) {
	query := `
		SELECT id, name, space_type, goal_type, units,
			   temperature_c, humidity_percent, speed_of_sound,
			   created_at, updated_at
		FROM projects 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*domain.Project

	for rows.Next() {
		var project domain.Project
		var spaceType, goalType, units string

		err := rows.Scan(
			&project.ID, &project.Name, &spaceType, &goalType, &units,
			&project.TemperatureC, &project.HumidityPercent, &project.SpeedOfSound,
			&project.CreatedAt, &project.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		project.SpaceType = domain.SpaceType(spaceType)
		project.GoalType = domain.GoalType(goalType)
		project.Units = domain.Units(units)

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	query := `
		UPDATE projects 
		SET name = ?, space_type = ?, goal_type = ?, units = ?,
			temperature_c = ?, humidity_percent = ?, speed_of_sound = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		project.Name, project.SpaceType, project.GoalType, project.Units,
		project.TemperatureC, project.HumidityPercent, project.SpeedOfSound,
		project.UpdatedAt, project.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", project.ID)
	}

	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", id)
	}

	return nil
}
