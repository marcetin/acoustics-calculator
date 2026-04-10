package repo

import (
	"context"
	"database/sql"
	"fmt"

	"acoustics-calculator/internal/domain"
)

type ConstraintRepository struct {
	db *sql.DB
}

func NewConstraintRepository(db *sql.DB) *ConstraintRepository {
	return &ConstraintRepository{db: db}
}

func (r *ConstraintRepository) Create(ctx context.Context, constraint *domain.ConstraintSet) error {
	query := `
		INSERT INTO constraint_sets (
			id, project_id, max_panel_depth_m, symmetry_required,
			allowed_surface_kinds_json, forbidden_surface_ids_json,
			preferred_treatment_mode, budget_class, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		constraint.ID, constraint.ProjectID, constraint.MaxPanelDepthM,
		constraint.SymmetryRequired, constraint.AllowedSurfaceKindsJSON,
		constraint.ForbiddenSurfaceIDsJSON, constraint.PreferredTreatmentMode,
		constraint.BudgetClass, constraint.CreatedAt, constraint.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create constraint: %w", err)
	}
	
	return nil
}

func (r *ConstraintRepository) GetByProjectID(ctx context.Context, projectID string) (*domain.ConstraintSet, error) {
	query := `
		SELECT id, project_id, max_panel_depth_m, symmetry_required,
			   allowed_surface_kinds_json, forbidden_surface_ids_json,
			   preferred_treatment_mode, budget_class, created_at, updated_at
		FROM constraint_sets 
		WHERE project_id = ?
	`
	
	var constraint domain.ConstraintSet
	var treatmentMode string
	var budgetClass string
	
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&constraint.ID, &constraint.ProjectID, &constraint.MaxPanelDepthM,
		&constraint.SymmetryRequired, &constraint.AllowedSurfaceKindsJSON,
		&constraint.ForbiddenSurfaceIDsJSON, &treatmentMode, &budgetClass,
		&constraint.CreatedAt, &constraint.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get constraint: %w", err)
	}
	
	constraint.PreferredTreatmentMode = domain.TreatmentMode(treatmentMode)
	constraint.BudgetClass = domain.BudgetClass(budgetClass)
	
	return &constraint, nil
}

func (r *ConstraintRepository) GetByID(ctx context.Context, id string) (*domain.ConstraintSet, error) {
	query := `
		SELECT id, project_id, max_panel_depth_m, symmetry_required,
			   allowed_surface_kinds_json, forbidden_surface_ids_json,
			   preferred_treatment_mode, budget_class, created_at, updated_at
		FROM constraint_sets 
		WHERE id = ?
	`
	
	var constraint domain.ConstraintSet
	var treatmentMode string
	var budgetClass string
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&constraint.ID, &constraint.ProjectID, &constraint.MaxPanelDepthM,
		&constraint.SymmetryRequired, &constraint.AllowedSurfaceKindsJSON,
		&constraint.ForbiddenSurfaceIDsJSON, &treatmentMode, &budgetClass,
		&constraint.CreatedAt, &constraint.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("constraint not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get constraint: %w", err)
	}
	
	constraint.PreferredTreatmentMode = domain.TreatmentMode(treatmentMode)
	constraint.BudgetClass = domain.BudgetClass(budgetClass)
	
	return &constraint, nil
}

func (r *ConstraintRepository) Update(ctx context.Context, constraint *domain.ConstraintSet) error {
	query := `
		UPDATE constraint_sets 
		SET max_panel_depth_m = ?, symmetry_required = ?,
		    allowed_surface_kinds_json = ?, forbidden_surface_ids_json = ?,
		    preferred_treatment_mode = ?, budget_class = ?, updated_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		constraint.MaxPanelDepthM, constraint.SymmetryRequired,
		constraint.AllowedSurfaceKindsJSON, constraint.ForbiddenSurfaceIDsJSON,
		constraint.PreferredTreatmentMode, constraint.BudgetClass,
		constraint.UpdatedAt, constraint.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update constraint: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("constraint not found: %s", constraint.ID)
	}
	
	return nil
}

func (r *ConstraintRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM constraint_sets WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete constraint: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("constraint not found: %s", id)
	}
	
	return nil
}

func (r *ConstraintRepository) DeleteByProjectID(ctx context.Context, projectID string) error {
	query := `DELETE FROM constraint_sets WHERE project_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete constraint by project: %w", err)
	}
	
	return nil
}
