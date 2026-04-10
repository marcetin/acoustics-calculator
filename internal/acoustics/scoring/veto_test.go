package scoring

import (
	"testing"

	"acoustics-calculator/internal/acoustics/diffuser"
	"acoustics-calculator/internal/domain"
)

func TestVetoChecker_NotMountable(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: false},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: []domain.ReflectionPath{},
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	candidate := Candidate{
		SurfaceID: "wall1",
		Surface:   &surfaces[0],
	}

	vetoChecker := NewVetoChecker(domain.SpaceTypeStudio, nil)
	result := vetoChecker.Apply(&candidate, 0.5, ctx)

	if result.Decision != domain.PlacementDecisionReject {
		t.Errorf("Non-mountable surface should be REJECT, got %s", result.Decision)
	}
}

func TestVetoChecker_PatchTooSmall(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: true},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: []domain.ReflectionPath{},
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	patch := &domain.BoundingBox{
		MinX: -0.1, MaxX: 0.1,
		MinY: -0.1, MaxY: 0.1,
	}

	candidate := Candidate{
		SurfaceID: "wall1",
		Surface:   &surfaces[0],
		Patch:     patch,
	}

	vetoChecker := NewVetoChecker(domain.SpaceTypeStudio, nil)
	result := vetoChecker.Apply(&candidate, 0.5, ctx)

	if result.Decision != domain.PlacementDecisionReject {
		t.Errorf("Too small patch should be REJECT, got %s", result.Decision)
	}
}

func TestVetoChecker_NearfieldStudio(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: true},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: []domain.ReflectionPath{},
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	patch := &domain.BoundingBox{
		MinX: -0.5, MaxX: 0.5,
		MinY: -0.5, MaxY: 0.5,
	}

	candidate := Candidate{
		SurfaceID:        "wall1",
		Surface:          &surfaces[0],
		Patch:            patch,
		AvgArrivalTimeMs: 10.0,
	}

	vetoChecker := NewVetoChecker(domain.SpaceTypeStudio, nil)
	result := vetoChecker.Apply(&candidate, 0.5, ctx)

	if result.Decision != domain.PlacementDecisionAbsorberRecommended {
		t.Errorf("Nearfield in STUDIO should be ABSORBER_RECOMMENDED, got %s", result.Decision)
	}
}

func TestVetoChecker_FirstReflectionStudio(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: true},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: []domain.ReflectionPath{},
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	patch := &domain.BoundingBox{
		MinX: -0.5, MaxX: 0.5,
		MinY: -0.5, MaxY: 0.5,
	}

	candidate := Candidate{
		SurfaceID:        "wall1",
		Surface:          &surfaces[0],
		Patch:            patch,
		AvgArrivalTimeMs: 15.0,
		ReflectionCount:  3,
	}

	vetoChecker := NewVetoChecker(domain.SpaceTypeStudio, nil)
	result := vetoChecker.Apply(&candidate, 0.5, ctx)

	if result.Decision != domain.PlacementDecisionAbsorberRecommended {
		t.Errorf("First reflection sidewall in STUDIO should be ABSORBER_RECOMMENDED, got %s", result.Decision)
	}
}

func TestVetoChecker_LowScore(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: true},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: []domain.ReflectionPath{},
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	patch := &domain.BoundingBox{
		MinX: -0.5, MaxX: 0.5,
		MinY: -0.5, MaxY: 0.5,
	}

	candidate := Candidate{
		SurfaceID:        "wall1",
		Surface:          &surfaces[0],
		Patch:            patch,
		AvgArrivalTimeMs: 30.0,
	}

	vetoChecker := NewVetoChecker(domain.SpaceTypeHifi, nil)
	result := vetoChecker.Apply(&candidate, 0.2, ctx)

	if result.Decision != domain.PlacementDecisionReject {
		t.Errorf("Low score should be REJECT, got %s", result.Decision)
	}
}
