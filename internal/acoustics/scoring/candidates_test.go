package scoring

import (
	"testing"

	"acoustics-calculator/internal/acoustics/diffuser"
	"acoustics-calculator/internal/domain"
)

func TestCandidateExtractor_Extract(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: true},
	}

	reflections := []domain.ReflectionPath{
		{SurfaceSequence: []string{"wall1"}, EnergyEstimate: 0.5, ArrivalTimeMs: 20.0, SourceID: "s1", ReceiverID: "r1"},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: reflections,
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	extractor := NewCandidateExtractor()
	candidates := extractor.Extract(ctx)

	if len(candidates) == 0 {
		t.Error("Extract() should return candidates")
	}
}

func TestCandidateExtractor_Deterministic(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: true},
	}

	reflections := []domain.ReflectionPath{
		{SurfaceSequence: []string{"wall1"}, EnergyEstimate: 0.5, ArrivalTimeMs: 20.0, SourceID: "s1", ReceiverID: "r1"},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: reflections,
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	extractor := NewCandidateExtractor()
	candidates1 := extractor.Extract(ctx)
	candidates2 := extractor.Extract(ctx)

	if len(candidates1) != len(candidates2) {
		t.Errorf("Deterministic check failed: different candidate counts")
	}
}

func TestCandidateExtractor_NonMountableFiltered(t *testing.T) {
	surfaces := []domain.Surface{
		{ID: "wall1", Kind: domain.SurfaceKindWall, AreaM2: 10.0, IsMountable: false},
	}

	reflections := []domain.ReflectionPath{
		{SurfaceSequence: []string{"wall1"}, EnergyEstimate: 0.5, ArrivalTimeMs: 20.0, SourceID: "s1", ReceiverID: "r1"},
	}

	metrics := &domain.AnalysisMetrics{
		Reflections: reflections,
	}

	ctx := &CandidateContext{
		Surfaces:        surfaces,
		LatestAnalysis:  metrics,
		DiffuserCatalog: diffuser.NewCatalog([]*domain.DiffuserType{}),
	}

	extractor := NewCandidateExtractor()
	candidates := extractor.Extract(ctx)

	if len(candidates) > 0 {
		t.Error("Non-mountable surfaces should be filtered")
	}
}
