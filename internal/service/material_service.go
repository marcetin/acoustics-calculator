package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type MaterialService struct {
	materialRepo *repo.MaterialRepository
	surfaceRepo  *repo.SurfaceRepository
}

func NewMaterialService(materialRepo *repo.MaterialRepository, surfaceRepo *repo.SurfaceRepository) *MaterialService {
	return &MaterialService{
		materialRepo: materialRepo,
		surfaceRepo:  surfaceRepo,
	}
}

func (s *MaterialService) ListPresets(ctx context.Context) ([]*domain.Material, error) {
	return s.materialRepo.ListPresets(ctx)
}

func (s *MaterialService) ListProjectMaterials(ctx context.Context, projectID string) ([]*domain.Material, error) {
	return s.materialRepo.ListByProjectID(ctx, projectID)
}

func (s *MaterialService) AssignMaterialToSurface(ctx context.Context, surfaceID string, materialID *string) error {
	surface, err := s.surfaceRepo.GetByID(ctx, surfaceID)
	if err != nil {
		return fmt.Errorf("failed to get surface: %w", err)
	}

	if materialID != nil {
		_, err = s.materialRepo.GetByID(ctx, *materialID)
		if err != nil {
			return fmt.Errorf("material not found: %w", err)
		}
	}

	surface.MaterialID = materialID
	surface.Update(domain.UpdateSurfaceInput{
		MaterialID: materialID,
	})

	if err := s.surfaceRepo.Update(ctx, surface); err != nil {
		return fmt.Errorf("failed to update surface: %w", err)
	}

	return nil
}

func (s *MaterialService) ListProjectSurfaces(ctx context.Context, projectID string) ([]*domain.Surface, error) {
	return s.surfaceRepo.ListByProjectID(ctx, projectID)
}

func (s *MaterialService) CreateCustomMaterial(ctx context.Context, input domain.CreateMaterialInput) (*domain.Material, error) {
	if err := s.validateMaterialInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	input.IsPreset = false
	material, err := domain.NewMaterial(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create material: %w", err)
	}

	if err := s.materialRepo.Create(ctx, material); err != nil {
		return nil, fmt.Errorf("failed to save material: %w", err)
	}

	return material, nil
}

func (s *MaterialService) validateMaterialInput(input domain.CreateMaterialInput) error {
	var errors []string

	if input.Name == "" {
		errors = append(errors, "material name is required")
	}

	if input.Category == "" {
		errors = append(errors, "material category is required")
	}

	if len(input.AbsorptionBands) == 0 {
		errors = append(errors, "absorption bands are required")
	}

	for _, band := range input.AbsorptionBands {
		if band.Value < 0 || band.Value > 1.0 {
			errors = append(errors, "absorption value must be between 0 and 1")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}
