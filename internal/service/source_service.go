package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type SourceService struct {
	sourceRepo   *repo.SourceRepository
	geometryRepo *repo.GeometryRepository
}

func NewSourceService(sourceRepo *repo.SourceRepository, geometryRepo *repo.GeometryRepository) *SourceService {
	return &SourceService{
		sourceRepo:   sourceRepo,
		geometryRepo: geometryRepo,
	}
}

func (s *SourceService) ListByProject(ctx context.Context, projectID string) ([]*domain.Source, error) {
	return s.sourceRepo.ListByProjectID(ctx, projectID)
}

func (s *SourceService) GetByID(ctx context.Context, id string) (*domain.Source, error) {
	return s.sourceRepo.GetByID(ctx, id)
}

func (s *SourceService) Create(ctx context.Context, input domain.CreateSourceInput) (*domain.Source, error) {
	if err := s.validateSourceInput(ctx, input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	source := domain.NewSource(input)

	if err := s.sourceRepo.Create(ctx, source); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	return source, nil
}

func (s *SourceService) Update(ctx context.Context, id string, input domain.UpdateSourceInput) (*domain.Source, error) {
	source, err := s.sourceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if err := s.validateSourceUpdate(ctx, source.ProjectID, input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	source.Update(input)

	if err := s.sourceRepo.Update(ctx, source); err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	return source, nil
}

func (s *SourceService) Delete(ctx context.Context, id string) error {
	return s.sourceRepo.Delete(ctx, id)
}

func (s *SourceService) validateSourceInput(ctx context.Context, input domain.CreateSourceInput) error {
	var errors []string

	if input.ProjectID == "" {
		errors = append(errors, "project ID is required")
	}

	if input.Name == "" {
		errors = append(errors, "source name is required")
	}

	if input.PositionZ < 0 {
		errors = append(errors, "Z position cannot be negative")
	}

	if err := s.validateWithinGeometry(ctx, input.ProjectID, input.PositionX, input.PositionY, input.PositionZ); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}

func (s *SourceService) validateSourceUpdate(ctx context.Context, projectID string, input domain.UpdateSourceInput) error {
	var errors []string

	if input.PositionX != nil && input.PositionY != nil && input.PositionZ != nil {
		if *input.PositionZ < 0 {
			errors = append(errors, "Z position cannot be negative")
		}

		if err := s.validateWithinGeometry(ctx, projectID, *input.PositionX, *input.PositionY, *input.PositionZ); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}

func (s *SourceService) validateWithinGeometry(ctx context.Context, projectID string, x, y, z float64) error {
	geometry, err := s.geometryRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get geometry: %w", err)
	}

	if geometry == nil {
		return fmt.Errorf("room geometry must be defined before adding sources")
	}

	if !geometry.IsShoebox() {
		return fmt.Errorf("position validation not yet implemented for polyhedral geometry")
	}

	if x < 0 || x > geometry.Width {
		return fmt.Errorf("X position must be between 0 and room width")
	}

	if y < 0 || y > geometry.Length {
		return fmt.Errorf("Y position must be between 0 and room length")
	}

	if z < 0 || z > geometry.Height {
		return fmt.Errorf("Z position must be between 0 and room height")
	}

	return nil
}
