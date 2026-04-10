package acoustics

import (
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestRTEstimator_CalculateRT(t *testing.T) {
	surfaces := []domain.Surface{
		{AreaM2: 30.0},
		{AreaM2: 30.0},
		{AreaM2: 18.0},
		{AreaM2: 18.0},
		{AreaM2: 15.0},
		{AreaM2: 15.0},
	}

	materials := map[string]*domain.Material{}

	estimator := NewRTEstimator(90.0, 343)
	rtResults, warnings := estimator.CalculateRT(surfaces, materials)

	if len(rtResults.Bands) != len(OctaveBands) {
		t.Errorf("CalculateRT() returned %d bands, want %d", len(rtResults.Bands), len(OctaveBands))
	}

	if len(warnings) == 0 {
		t.Error("CalculateRT() should produce warnings for missing materials")
	}

	for _, band := range rtResults.Bands {
		if band.SabineT60 < 0 {
			t.Errorf("SabineT60 should be non-negative, got %f", band.SabineT60)
		}
		if band.EyringT60 < 0 {
			t.Errorf("EyringT60 should be non-negative, got %f", band.EyringT60)
		}
		if band.MeanAbsorption < 0 {
			t.Errorf("MeanAbsorption should be non-negative, got %f", band.MeanAbsorption)
		}
	}
}

func TestRTEstimator_NoNaNInf(t *testing.T) {
	surfaces := []domain.Surface{
		{AreaM2: 10.0},
		{AreaM2: 10.0},
		{AreaM2: 10.0},
		{AreaM2: 10.0},
		{AreaM2: 10.0},
		{AreaM2: 10.0},
	}

	materials := map[string]*domain.Material{}

	estimator := NewRTEstimator(50.0, 343)
	rtResults, _ := estimator.CalculateRT(surfaces, materials)

	for _, band := range rtResults.Bands {
		if isInvalid(band.SabineT60) {
			t.Errorf("SabineT60 is invalid: %f", band.SabineT60)
		}
		if isInvalid(band.EyringT60) {
			t.Errorf("EyringT60 is invalid: %f", band.EyringT60)
		}
		if isInvalid(band.MeanAbsorption) {
			t.Errorf("MeanAbsorption is invalid: %f", band.MeanAbsorption)
		}
	}
}

func isInvalid(f float64) bool {
	return f != f || f > 1e300 || f < -1e300
}
