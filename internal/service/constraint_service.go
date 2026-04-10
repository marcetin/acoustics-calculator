package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type ConstraintService struct {
	constraintRepo *repo.ConstraintRepository
}

func NewConstraintService(constraintRepo *repo.ConstraintRepository) *ConstraintService {
	return &ConstraintService{
		constraintRepo: constraintRepo,
	}
}

func (s *ConstraintService) GetOrCreateDefault(ctx context.Context, projectID string) (*domain.ConstraintSet, error) {
	constraint, err := s.constraintRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get constraint: %w", err)
	}

	if constraint != nil {
		return constraint, nil
	}

	defaultInput := domain.GetDefaultConstraints(projectID)
	newConstraint, err := domain.NewConstraintSet(defaultInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create default constraint: %w", err)
	}

	if err := s.constraintRepo.Create(ctx, newConstraint); err != nil {
		return nil, fmt.Errorf("failed to save default constraint: %w", err)
	}

	return newConstraint, nil
}

func (s *ConstraintService) Save(ctx context.Context, projectID string, input domain.UpdateConstraintInput) (*domain.ConstraintSet, error) {
	if err := s.validateConstraintInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	constraint, err := s.GetOrCreateDefault(ctx, projectID)
	if err != nil {
		return nil, err
	}

	constraint.Update(input)

	if err := s.constraintRepo.Update(ctx, constraint); err != nil {
		return nil, fmt.Errorf("failed to update constraint: %w", err)
	}

	return constraint, nil
}

func (s *ConstraintService) validateConstraintInput(input domain.UpdateConstraintInput) error {
	var errors []string

	if input.MaxPanelDepthM != nil {
		if *input.MaxPanelDepthM <= 0 {
			errors = append(errors, "max panel depth must be greater than 0")
		}
		if *input.MaxPanelDepthM > 2.0 {
			errors = append(errors, "max panel depth must be less than 2 meters")
		}
	}

	if input.AllowedSurfaceKinds != nil {
		for _, kind := range *input.AllowedSurfaceKinds {
			valid := false
			for _, validKind := range domain.ValidSurfaceKinds() {
				if kind == validKind {
					valid = true
					break
				}
			}
			if !valid {
				errors = append(errors, fmt.Sprintf("invalid surface kind: %s", kind))
			}
		}
	}

	if input.PreferredTreatmentMode != nil {
		valid := false
		for _, validMode := range domain.ValidTreatmentModes() {
			if *input.PreferredTreatmentMode == validMode {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, "invalid treatment mode")
		}
	}

	if input.BudgetClass != nil {
		valid := false
		for _, validClass := range domain.ValidBudgetClasses() {
			if *input.BudgetClass == validClass {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, "invalid budget class")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}
