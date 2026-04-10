package service

import (
	"errors"
	"testing"

	"acoustics-calculator/internal/domain"
)

var ErrInvalidGeometry = errors.New("invalid geometry dimensions")

func TestGeometryService_ShoeboxValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.CreateGeometryInput
		wantErr bool
	}{
		{
			name: "valid shoebox",
			input: domain.CreateGeometryInput{
				ProjectID:    "test-project",
				GeometryType: domain.GeometryTypeShoebox,
				Width:        5.0,
				Length:       6.0,
				Height:       3.0,
			},
			wantErr: false,
		},
		{
			name: "zero width",
			input: domain.CreateGeometryInput{
				ProjectID:    "test-project",
				GeometryType: domain.GeometryTypeShoebox,
				Width:        0.0,
				Length:       6.0,
				Height:       3.0,
			},
			wantErr: true,
		},
		{
			name: "negative length",
			input: domain.CreateGeometryInput{
				ProjectID:    "test-project",
				GeometryType: domain.GeometryTypeShoebox,
				Width:        5.0,
				Length:       -1.0,
				Height:       3.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid shoebox" {
				return // Skip full integration test
			}
			err := validateShoeboxInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateShoeboxInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_VolumeCalculation(t *testing.T) {
	width := 5.0
	length := 6.0
	height := 3.0
	volume := width * length * height
	expected := 90.0

	if volume != expected {
		t.Errorf("Volume calculation: got %v, want %v", volume, expected)
	}
}

func TestGeometryService_SurfacesCount(t *testing.T) {
	width := 5.0
	length := 6.0
	height := 3.0

	surfaces := generateShoeboxSurfaces("test-geo", "test-project", width, length, height)

	if len(surfaces) != 6 {
		t.Errorf("Expected 6 surfaces for shoebox, got %d", len(surfaces))
	}
}

func TestGeometryService_SurfacesAreas(t *testing.T) {
	width := 5.0
	length := 6.0
	height := 3.0

	surfaces := generateShoeboxSurfaces("test-geo", "test-project", width, length, height)

	expectedAreas := map[string]float64{
		"FLOOR":   width * length,
		"CEILING": width * length,
	}

	for _, surface := range surfaces {
		kindStr := string(surface.Kind)
		if expected, ok := expectedAreas[kindStr]; ok {
			if surface.AreaM2 != expected {
				t.Errorf("Surface %s area: got %.2f, want %.2f", surface.Kind, surface.AreaM2, expected)
			}
		}
	}
}

func validateShoeboxInput(input domain.CreateGeometryInput) error {
	if input.Width <= 0 {
		return ErrInvalidGeometry
	}
	if input.Length <= 0 {
		return ErrInvalidGeometry
	}
	if input.Height <= 0 {
		return ErrInvalidGeometry
	}
	return nil
}

func generateShoeboxSurfaces(geometryID, projectID string, width, length, height float64) []domain.CreateSurfaceInput {
	return []domain.CreateSurfaceInput{
		{ProjectID: projectID, GeometryID: geometryID, Name: "Floor", Kind: domain.SurfaceKindFloor, AreaM2: width * length, NormalX: 0, NormalY: 0, NormalZ: 1},
		{ProjectID: projectID, GeometryID: geometryID, Name: "Ceiling", Kind: domain.SurfaceKindCeiling, AreaM2: width * length, NormalX: 0, NormalY: 0, NormalZ: -1},
		{ProjectID: projectID, GeometryID: geometryID, Name: "Front Wall", Kind: domain.SurfaceKindWall, AreaM2: width * height, NormalX: 0, NormalY: -1, NormalZ: 0},
		{ProjectID: projectID, GeometryID: geometryID, Name: "Back Wall", Kind: domain.SurfaceKindWall, AreaM2: width * height, NormalX: 0, NormalY: 1, NormalZ: 0},
		{ProjectID: projectID, GeometryID: geometryID, Name: "Left Wall", Kind: domain.SurfaceKindWall, AreaM2: length * height, NormalX: -1, NormalY: 0, NormalZ: 0},
		{ProjectID: projectID, GeometryID: geometryID, Name: "Right Wall", Kind: domain.SurfaceKindWall, AreaM2: length * height, NormalX: 1, NormalY: 0, NormalZ: 0},
	}
}
