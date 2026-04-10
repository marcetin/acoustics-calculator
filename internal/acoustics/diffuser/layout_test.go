package diffuser

import (
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestLayoutSynthesizer_SynthesizeLayout(t *testing.T) {
	patch := &domain.BoundingBox{
		MinX: -0.5, MaxX: 0.5,
		MinY: -0.5, MaxY: 0.5,
		MinZ: 0, MaxZ: 0,
	}

	diffuser := &domain.DiffuserType{
		PanelWidthM:  0.5,
		PanelHeightM: 0.5,
	}

	synthesizer := NewLayoutSynthesizer()
	result := synthesizer.SynthesizeLayout(patch, diffuser, "VERTICAL")

	if result == nil {
		t.Error("SynthesizeLayout should return a result")
	}

	if result.PanelCount == 0 {
		t.Error("PanelCount should be greater than 0")
	}

	if len(result.PanelCenters) != result.PanelCount {
		t.Error("PanelCenters length should match PanelCount")
	}
}

func TestLayoutSynthesizer_Deterministic(t *testing.T) {
	patch := &domain.BoundingBox{
		MinX: -0.5, MaxX: 0.5,
		MinY: -0.5, MaxY: 0.5,
		MinZ: 0, MaxZ: 0,
	}

	diffuser := &domain.DiffuserType{
		PanelWidthM:  0.5,
		PanelHeightM: 0.5,
	}

	synthesizer := NewLayoutSynthesizer()
	result1 := synthesizer.SynthesizeLayout(patch, diffuser, "VERTICAL")
	result2 := synthesizer.SynthesizeLayout(patch, diffuser, "VERTICAL")

	if result1.PanelCount != result2.PanelCount {
		t.Error("Deterministic check failed: different panel counts")
	}

	if len(result1.PanelCenters) != len(result2.PanelCenters) {
		t.Error("Deterministic check failed: different panel center counts")
	}
}

func TestLayoutSynthesizer_PanelCentersInBounds(t *testing.T) {
	patch := &domain.BoundingBox{
		MinX: -0.5, MaxX: 0.5,
		MinY: -0.5, MaxY: 0.5,
		MinZ: 0, MaxZ: 0,
	}

	diffuser := &domain.DiffuserType{
		PanelWidthM:  0.5,
		PanelHeightM: 0.5,
	}

	synthesizer := NewLayoutSynthesizer()
	result := synthesizer.SynthesizeLayout(patch, diffuser, "VERTICAL")

	for _, center := range result.PanelCenters {
		if center.X < patch.MinX || center.X > patch.MaxX {
			t.Errorf("Panel center X out of bounds: %f", center.X)
		}
		if center.Y < patch.MinY || center.Y > patch.MaxY {
			t.Errorf("Panel center Y out of bounds: %f", center.Y)
		}
	}
}
