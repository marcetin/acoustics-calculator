package acoustics

import (
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestImageSourceSolver_CalculateReflections(t *testing.T) {
	sources := []domain.Source{
		{ID: "s1", PositionX: 2.5, PositionY: 3.0, PositionZ: 1.2},
	}

	receivers := []domain.Receiver{
		{ID: "r1", PositionX: 2.5, PositionY: 4.0, PositionZ: 1.2},
	}

	surfaces := []domain.Surface{
		{ID: "floor", Kind: domain.SurfaceKindFloor, AreaM2: 30.0},
		{ID: "ceiling", Kind: domain.SurfaceKindCeiling, AreaM2: 30.0},
		{ID: "front", Kind: domain.SurfaceKindWall, AreaM2: 18.0},
		{ID: "back", Kind: domain.SurfaceKindWall, AreaM2: 18.0},
		{ID: "left", Kind: domain.SurfaceKindWall, AreaM2: 15.0},
		{ID: "right", Kind: domain.SurfaceKindWall, AreaM2: 15.0},
	}

	solver := NewImageSourceSolver(5.0, 6.0, 3.0, 343)
	reflections, warnings := solver.CalculateReflections(sources, receivers, surfaces)

	if len(reflections) == 0 {
		t.Error("CalculateReflections() returned no reflections")
	}

	if len(warnings) > 0 {
		t.Logf("Warnings: %v", warnings)
	}

	maxOrder := 0
	for _, ref := range reflections {
		if ref.Order > maxOrder {
			maxOrder = ref.Order
		}

		if ref.PathLengthM <= 0 {
			t.Errorf("PathLengthM should be positive, got %f", ref.PathLengthM)
		}
		if ref.ArrivalTimeMs <= 0 {
			t.Errorf("ArrivalTimeMs should be positive, got %f", ref.ArrivalTimeMs)
		}
		if ref.EnergyEstimate <= 0 {
			t.Errorf("EnergyEstimate should be positive, got %f", ref.EnergyEstimate)
		}
	}

	if maxOrder > 2 {
		t.Errorf("Max reflection order should be 2, got %d", maxOrder)
	}
}

func TestImageSourceSolver_Deterministic(t *testing.T) {
	sources := []domain.Source{
		{ID: "s1", PositionX: 2.5, PositionY: 3.0, PositionZ: 1.2},
	}

	receivers := []domain.Receiver{
		{ID: "r1", PositionX: 2.5, PositionY: 4.0, PositionZ: 1.2},
	}

	surfaces := []domain.Surface{
		{ID: "floor", Kind: domain.SurfaceKindFloor, AreaM2: 30.0},
		{ID: "ceiling", Kind: domain.SurfaceKindCeiling, AreaM2: 30.0},
		{ID: "front", Kind: domain.SurfaceKindWall, AreaM2: 18.0},
		{ID: "back", Kind: domain.SurfaceKindWall, AreaM2: 18.0},
		{ID: "left", Kind: domain.SurfaceKindWall, AreaM2: 15.0},
		{ID: "right", Kind: domain.SurfaceKindWall, AreaM2: 15.0},
	}

	solver := NewImageSourceSolver(5.0, 6.0, 3.0, 343)
	reflections1, _ := solver.CalculateReflections(sources, receivers, surfaces)
	reflections2, _ := solver.CalculateReflections(sources, receivers, surfaces)

	if len(reflections1) != len(reflections2) {
		t.Errorf("Deterministic check failed: different reflection counts %d vs %d", len(reflections1), len(reflections2))
	}
}

func TestImageSourceSolver_ArrivalTimeConsistency(t *testing.T) {
	sources := []domain.Source{
		{ID: "s1", PositionX: 2.5, PositionY: 3.0, PositionZ: 1.2},
	}

	receivers := []domain.Receiver{
		{ID: "r1", PositionX: 2.5, PositionY: 4.0, PositionZ: 1.2},
	}

	surfaces := []domain.Surface{
		{ID: "floor", Kind: domain.SurfaceKindFloor, AreaM2: 30.0},
		{ID: "ceiling", Kind: domain.SurfaceKindCeiling, AreaM2: 30.0},
		{ID: "front", Kind: domain.SurfaceKindWall, AreaM2: 18.0},
		{ID: "back", Kind: domain.SurfaceKindWall, AreaM2: 18.0},
		{ID: "left", Kind: domain.SurfaceKindWall, AreaM2: 15.0},
		{ID: "right", Kind: domain.SurfaceKindWall, AreaM2: 15.0},
	}

	c := 343.0
	solver := NewImageSourceSolver(5.0, 6.0, 3.0, c)
	reflections, _ := solver.CalculateReflections(sources, receivers, surfaces)

	for _, ref := range reflections {
		expectedTime := (ref.PathLengthM / c) * 1000
		diff := ref.ArrivalTimeMs - expectedTime
		if diff < -0.1 || diff > 0.1 {
			t.Errorf("Arrival time inconsistent: got %.2f ms, expected %.2f ms (diff %.2f)", ref.ArrivalTimeMs, expectedTime, diff)
		}
	}
}
