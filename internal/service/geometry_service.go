package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type GeometryService struct {
	geometryRepo *repo.GeometryRepository
	surfaceRepo  *repo.SurfaceRepository
}

func NewGeometryService(geometryRepo *repo.GeometryRepository, surfaceRepo *repo.SurfaceRepository) *GeometryService {
	return &GeometryService{
		geometryRepo: geometryRepo,
		surfaceRepo:  surfaceRepo,
	}
}

func (s *GeometryService) GetByProjectID(ctx context.Context, projectID string) (*domain.RoomGeometry, error) {
	return s.geometryRepo.GetByProjectID(ctx, projectID)
}

func (s *GeometryService) UpsertShoeboxGeometry(ctx context.Context, input domain.CreateGeometryInput) (*domain.RoomGeometry, error) {
	if err := s.validateGeometryInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existingGeometry, err := s.geometryRepo.GetByProjectID(ctx, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing geometry: %w", err)
	}

	var geometry *domain.RoomGeometry
	if existingGeometry != nil {
		if existingGeometry.GeometryType != domain.GeometryTypeShoebox {
			return nil, fmt.Errorf("project already has %s geometry, cannot switch to shoebox", existingGeometry.GeometryType)
		}

		updateInput := domain.UpdateGeometryInput{
			Width:  &input.Width,
			Length: &input.Length,
			Height: &input.Height,
		}
		existingGeometry.Update(updateInput)

		if err := s.geometryRepo.Update(ctx, existingGeometry); err != nil {
			return nil, fmt.Errorf("failed to update geometry: %w", err)
		}

		if err := s.regenerateShoeboxSurfaces(ctx, existingGeometry); err != nil {
			return nil, fmt.Errorf("failed to regenerate surfaces: %w", err)
		}

		geometry = existingGeometry
	} else {
		geometry = domain.NewRoomGeometry(input)

		if err := s.geometryRepo.Create(ctx, geometry); err != nil {
			return nil, fmt.Errorf("failed to create geometry: %w", err)
		}

		if err := s.generateShoeboxSurfaces(ctx, geometry); err != nil {
			return nil, fmt.Errorf("failed to generate surfaces: %w", err)
		}
	}

	return geometry, nil
}

func (s *GeometryService) validateGeometryInput(input domain.CreateGeometryInput) error {
	var errors []string

	if input.ProjectID == "" {
		errors = append(errors, "project ID is required")
	}

	if input.Width <= 0 {
		errors = append(errors, "width must be greater than 0")
	}
	if input.Width < 1.0 {
		errors = append(errors, "width must be at least 1 meter")
	}
	if input.Width > 100.0 {
		errors = append(errors, "width must be less than 100 meters")
	}

	if input.Length <= 0 {
		errors = append(errors, "length must be greater than 0")
	}
	if input.Length < 1.0 {
		errors = append(errors, "length must be at least 1 meter")
	}
	if input.Length > 100.0 {
		errors = append(errors, "length must be less than 100 meters")
	}

	if input.Height <= 0 {
		errors = append(errors, "height must be greater than 0")
	}
	if input.Height < 1.0 {
		errors = append(errors, "height must be at least 1 meter")
	}
	if input.Height > 30.0 {
		errors = append(errors, "height must be less than 30 meters")
	}

	if input.GeometryType != domain.GeometryTypeShoebox && input.GeometryType != domain.GeometryTypePolyhedral {
		errors = append(errors, "invalid geometry type")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}

func (s *GeometryService) generateShoeboxSurfaces(ctx context.Context, geometry *domain.RoomGeometry) error {
	surfaces := s.createShoeboxSurfaces(geometry)

	for _, surface := range surfaces {
		if err := s.surfaceRepo.Create(ctx, surface); err != nil {
			return fmt.Errorf("failed to create surface %s: %w", surface.Name, err)
		}
	}

	return nil
}

func (s *GeometryService) regenerateShoeboxSurfaces(ctx context.Context, geometry *domain.RoomGeometry) error {
	if err := s.surfaceRepo.DeleteByGeometryID(ctx, geometry.ID); err != nil {
		return fmt.Errorf("failed to delete old surfaces: %w", err)
	}

	return s.generateShoeboxSurfaces(ctx, geometry)
}

func (s *GeometryService) createShoeboxSurfaces(geometry *domain.RoomGeometry) []*domain.Surface {
	w, l, h := geometry.Width, geometry.Length, geometry.Height

	surfaceInputs := []domain.CreateSurfaceInput{
		{
			Name:      "Front Wall",
			Kind:      domain.SurfaceKindWall,
			Polygon3D: []float64{0, 0, 0, w, 0, 0, w, h, 0, 0, h, 0},
			NormalX:   0, NormalY: -1, NormalZ: 0,
			AreaM2:      w * h,
			IsMountable: true,
		},
		{
			Name:      "Rear Wall",
			Kind:      domain.SurfaceKindWall,
			Polygon3D: []float64{0, 0, l, w, 0, l, w, h, l, 0, h, l},
			NormalX:   0, NormalY: 1, NormalZ: 0,
			AreaM2:      w * h,
			IsMountable: true,
		},
		{
			Name:      "Left Wall",
			Kind:      domain.SurfaceKindWall,
			Polygon3D: []float64{0, 0, 0, 0, 0, l, 0, h, l, 0, h, 0},
			NormalX:   -1, NormalY: 0, NormalZ: 0,
			AreaM2:      l * h,
			IsMountable: true,
		},
		{
			Name:      "Right Wall",
			Kind:      domain.SurfaceKindWall,
			Polygon3D: []float64{w, 0, 0, w, 0, l, w, h, l, w, h, 0},
			NormalX:   1, NormalY: 0, NormalZ: 0,
			AreaM2:      l * h,
			IsMountable: true,
		},
		{
			Name:      "Floor",
			Kind:      domain.SurfaceKindFloor,
			Polygon3D: []float64{0, 0, 0, w, 0, 0, w, 0, l, 0, 0, l},
			NormalX:   0, NormalY: 0, NormalZ: -1,
			AreaM2:      w * l,
			IsMountable: false,
		},
		{
			Name:      "Ceiling",
			Kind:      domain.SurfaceKindCeiling,
			Polygon3D: []float64{0, h, 0, w, h, 0, w, h, l, 0, h, l},
			NormalX:   0, NormalY: 0, NormalZ: 1,
			AreaM2:      w * l,
			IsMountable: true,
		},
	}

	var surfaces []*domain.Surface

	for _, input := range surfaceInputs {
		input.ProjectID = geometry.ProjectID
		input.GeometryID = geometry.ID

		surface, err := domain.NewSurface(input)
		if err != nil {
			continue
		}
		surfaces = append(surfaces, surface)
	}

	return surfaces
}
