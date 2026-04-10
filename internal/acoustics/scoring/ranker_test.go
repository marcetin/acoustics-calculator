package scoring

import (
	"testing"

	"acoustics-calculator/internal/acoustics/diffuser"
	"acoustics-calculator/internal/domain"
)

func TestRanker_Score(t *testing.T) {
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

	candidate := Candidate{
		SurfaceID:        "wall1",
		Surface:          &surfaces[0],
		ReflectionCount:  5,
		TotalEnergy:      0.05,
		AvgArrivalTimeMs: 25.0,
		Patch:            &domain.BoundingBox{MinX: -0.5, MaxX: 0.5, MinY: -0.5, MaxY: 0.5},
		TargetBands:      []float64{500, 1000},
	}

	ranker := NewRanker(domain.SpaceTypeStudio)
	score, _, _ := ranker.Score(&candidate, ctx)

	if score < 0 || score > 1 {
		t.Errorf("Score should be between 0 and 1, got %f", score)
	}
}

func TestRanker_SymmetryAffectsStudio(t *testing.T) {
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

	candidate := Candidate{
		SurfaceID:        "wall1",
		Surface:          &surfaces[0],
		ReflectionCount:  5,
		TotalEnergy:      0.05,
		AvgArrivalTimeMs: 25.0,
		Patch:            &domain.BoundingBox{MinX: -0.5, MaxX: 0.5, MinY: -0.5, MaxY: 0.5},
		TargetBands:      []float64{500, 1000},
	}

	studioRanker := NewRanker(domain.SpaceTypeStudio)
	hifiRanker := NewRanker(domain.SpaceTypeHifi)

	_, studioComponents, _ := studioRanker.Score(&candidate, ctx)
	_, hifiComponents, _ := hifiRanker.Score(&candidate, ctx)

	if studioComponents.SymmetryScore == hifiComponents.SymmetryScore {
		t.Error("Studio profile should weight symmetry differently")
	}
}

func TestRanker_NearfieldPenalty(t *testing.T) {
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

	nearfieldCandidate := Candidate{
		SurfaceID:        "wall1",
		Surface:          &surfaces[0],
		ReflectionCount:  5,
		TotalEnergy:      0.05,
		AvgArrivalTimeMs: 10.0,
		Patch:            &domain.BoundingBox{MinX: -0.5, MaxX: 0.5, MinY: -0.5, MaxY: 0.5},
		TargetBands:      []float64{500, 1000},
	}

	ranker := NewRanker(domain.SpaceTypeStudio)
	_, _, reasons := ranker.Score(&nearfieldCandidate, ctx)

	hasNearfieldPenalty := false
	for _, r := range reasons {
		if r.Code == "NEARFIELD_PENALTY" {
			hasNearfieldPenalty = true
			break
		}
	}

	if !hasNearfieldPenalty {
		t.Error("Nearfield candidate should have nearfield penalty reason")
	}
}
