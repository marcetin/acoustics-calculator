package scoring

import (
	"math"

	"acoustics-calculator/internal/acoustics/diffuser"
	"acoustics-calculator/internal/domain"
)

type CandidateExtractor struct {
	layoutSynthesizer *diffuser.LayoutSynthesizer
}

func NewCandidateExtractor() *CandidateExtractor {
	return &CandidateExtractor{
		layoutSynthesizer: diffuser.NewLayoutSynthesizer(),
	}
}

type CandidateContext struct {
	Project       *domain.Project
	Geometry      *domain.RoomGeometry
	Constraints   *domain.ConstraintSet
	LatestAnalysis *domain.AnalysisMetrics
	Surfaces      []domain.Surface
	DiffuserCatalog *diffuser.Catalog
}

type Candidate struct {
	SurfaceID           string
	Surface             *domain.Surface
	ReflectionCount     int
	TotalEnergy         float64
	AvgArrivalTimeMs    float64
	SourceReceiverPairs []string
	TargetBands         []float64
	Patch               *domain.BoundingBox
	Center              domain.Vec3
	Orientation         string
}

func (e *CandidateExtractor) Extract(ctx *CandidateContext) []Candidate {
	surfaceReflections := e.aggregateReflectionsBySurface(ctx.LatestAnalysis.Reflections, ctx.Surfaces)

	var candidates []Candidate
	for surfaceID, reflData := range surfaceReflections {
		if reflData.Count == 0 {
			continue
		}

		surface := e.findSurface(ctx.Surfaces, surfaceID)
		if surface == nil {
			continue
		}

		if !surface.IsMountable {
			continue
		}

		patch := e.createPatchFromReflections(reflData, surface)
		if patch == nil {
			continue
		}

		center := domain.Vec3{
			X: (patch.MinX + patch.MaxX) / 2,
			Y: (patch.MinY + patch.MaxY) / 2,
			Z: (patch.MinZ + patch.MaxZ) / 2,
		}

		orientation := e.determineOrientation(surface.Kind)

		targetBands := e.inferTargetBands(reflData)

		candidate := Candidate{
			SurfaceID:           surfaceID,
			Surface:             surface,
			ReflectionCount:     reflData.Count,
			TotalEnergy:         reflData.TotalEnergy,
			AvgArrivalTimeMs:    reflData.AvgArrivalTime,
			SourceReceiverPairs: reflData.SourceReceiverPairs,
			TargetBands:         targetBands,
			Patch:               patch,
			Center:              center,
			Orientation:         orientation,
		}

		candidates = append(candidates, candidate)
	}

	return candidates
}

type ReflectionData struct {
	Count              int
	TotalEnergy        float64
	AvgArrivalTime     float64
	SourceReceiverPairs []string
}

func (e *CandidateExtractor) aggregateReflectionsBySurface(reflections []domain.ReflectionPath, surfaces []domain.Surface) map[string]*ReflectionData {
	surfaceMap := make(map[string]*domain.Surface)
	for _, s := range surfaces {
		surfaceMap[s.ID] = &s
	}

	data := make(map[string]*ReflectionData)

	for _, ref := range reflections {
		if len(ref.SurfaceSequence) == 0 {
			continue
		}

		surfaceID := ref.SurfaceSequence[0]
		if _, exists := data[surfaceID]; !exists {
			data[surfaceID] = &ReflectionData{
				SourceReceiverPairs: []string{},
			}
		}

		d := data[surfaceID]
		d.Count++
		d.TotalEnergy += ref.EnergyEstimate
		d.AvgArrivalTime += ref.ArrivalTimeMs

		pair := ref.SourceID + "-" + ref.ReceiverID
		d.SourceReceiverPairs = append(d.SourceReceiverPairs, pair)
	}

	for _, d := range data {
		if d.Count > 0 {
			d.AvgArrivalTime = d.AvgArrivalTime / float64(d.Count)
		}
	}

	return data
}

func (e *CandidateExtractor) findSurface(surfaces []domain.Surface, id string) *domain.Surface {
	for i := range surfaces {
		if surfaces[i].ID == id {
			return &surfaces[i]
		}
	}
	return nil
}

func (e *CandidateExtractor) createPatchFromReflections(reflData *ReflectionData, surface *domain.Surface) *domain.BoundingBox {
	minPatchSize := 0.6

	width := math.Sqrt(surface.AreaM2)
	if width < minPatchSize {
		return nil
	}

	patchWidth := width
	patchHeight := width

	if patchWidth < minPatchSize {
		patchWidth = minPatchSize
	}
	if patchHeight < minPatchSize {
		patchHeight = minPatchSize
	}

	return &domain.BoundingBox{
		MinX: -patchWidth / 2,
		MaxX: patchWidth / 2,
		MinY: -patchHeight / 2,
		MaxY: patchHeight / 2,
		MinZ: 0,
		MaxZ: 0,
	}
}

func (e *CandidateExtractor) determineOrientation(kind domain.SurfaceKind) string {
	switch kind {
	case domain.SurfaceKindFloor:
		return "HORIZONTAL"
	case domain.SurfaceKindCeiling:
		return "HORIZONTAL"
	default:
		return "VERTICAL"
	}
}

func (e *CandidateExtractor) inferTargetBands(reflData *ReflectionData) []float64 {
	return []float64{500, 1000, 2000}
}
