package diffuser

import (
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestCatalog_FindCompatible(t *testing.T) {
	diffusers := []*domain.DiffuserType{
		{ID: "qrd1d", Category: domain.DiffuserCategoryQRD1D, PanelWidthM: 0.6, PanelHeightM: 0.6, PanelDepthM: 0.1, MinEffectiveHz: 500, MaxEffectiveHz: 4000, ScatterAxis: domain.ScatterAxisBoth},
		{ID: "qrd2d", Category: domain.DiffuserCategoryQRD2D, PanelWidthM: 1.2, PanelHeightM: 1.2, PanelDepthM: 0.15, MinEffectiveHz: 600, MaxEffectiveHz: 5000, ScatterAxis: domain.ScatterAxisBoth},
	}

	catalog := NewCatalog(diffusers)

	targetBands := []float64{500, 1000}
	maxDepth := 0.2
	patchArea := 1.0

	result := catalog.FindCompatible(targetBands, maxDepth, patchArea, "VERTICAL")

	if result == nil {
		t.Error("FindCompatible should return a diffuser")
	}
}

func TestCatalog_DepthLimit(t *testing.T) {
	diffusers := []*domain.DiffuserType{
		{ID: "deep", Category: domain.DiffuserCategoryQRD2D, PanelWidthM: 1.0, PanelHeightM: 1.0, PanelDepthM: 0.5, MinEffectiveHz: 500, MaxEffectiveHz: 4000, ScatterAxis: domain.ScatterAxisBoth},
	}

	catalog := NewCatalog(diffusers)

	targetBands := []float64{500, 1000}
	maxDepth := 0.1
	patchArea := 1.0

	result := catalog.FindCompatible(targetBands, maxDepth, patchArea, "VERTICAL")

	if result != nil {
		t.Error("FindCompatible should return nil when depth limit exceeded")
	}
}

func TestCatalog_PatchAreaLimit(t *testing.T) {
	diffusers := []*domain.DiffuserType{
		{ID: "large", Category: domain.DiffuserCategoryQRD2D, PanelWidthM: 2.0, PanelHeightM: 2.0, PanelDepthM: 0.1, MinEffectiveHz: 500, MaxEffectiveHz: 4000, ScatterAxis: domain.ScatterAxisBoth},
	}

	catalog := NewCatalog(diffusers)

	targetBands := []float64{500, 1000}
	maxDepth := 0.2
	patchArea := 0.5

	result := catalog.FindCompatible(targetBands, maxDepth, patchArea, "VERTICAL")

	if result != nil {
		t.Error("FindCompatible should return nil when panel area exceeds patch area")
	}
}
