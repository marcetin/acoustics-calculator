package scoring

import (
	"math"

	"acoustics-calculator/internal/domain"
)

type Ranker struct {
	projectSpaceType domain.SpaceType
}

func NewRanker(spaceType domain.SpaceType) *Ranker {
	return &Ranker{projectSpaceType: spaceType}
}

type ScoreComponents struct {
	EnergyScore      float64
	TimeScore        float64
	ReceiverScore    float64
	MountScore       float64
	SymmetryScore    float64
	DistanceScore    float64
	BandMatchScore   float64
	PathSupportScore float64
	ModalPenalty     float64
	NearfieldPenalty float64
}

func (r *Ranker) Score(candidate *Candidate, ctx *CandidateContext) (float64, *ScoreComponents, []domain.CandidateReason) {
	var reasons []domain.CandidateReason

	energyScore := r.normalizeEnergy(candidate.TotalEnergy, candidate.ReflectionCount)
	timeScore := r.normalizeArrivalTime(candidate.AvgArrivalTimeMs)
	receiverScore := r.calculateReceiverScore(candidate)
	mountScore := 1.0
	symmetryScore := r.calculateSymmetryScore(candidate, ctx)
	distanceScore := r.calculateDistanceScore(candidate, ctx)
	bandMatchScore := r.calculateBandMatchScore(candidate)
	pathSupportScore := r.calculatePathSupportScore(candidate)
	modalPenalty := r.calculateModalPenalty(candidate, ctx)
	nearfieldPenalty := r.calculateNearfieldPenalty(candidate, ctx)

	components := &ScoreComponents{
		EnergyScore:      energyScore,
		TimeScore:        timeScore,
		ReceiverScore:    receiverScore,
		MountScore:       mountScore,
		SymmetryScore:    symmetryScore,
		DistanceScore:    distanceScore,
		BandMatchScore:   bandMatchScore,
		PathSupportScore: pathSupportScore,
		ModalPenalty:     modalPenalty,
		NearfieldPenalty: nearfieldPenalty,
	}

	if symmetryScore > 0.8 {
		reasons = append(reasons, domain.CandidateReason{
			Code:    "SYMMETRY_PRESERVED",
			Message: "Symmetry preserved in placement",
		})
	}

	if nearfieldPenalty > 0.3 {
		reasons = append(reasons, domain.CandidateReason{
			Code:    "NEARFIELD_PENALTY",
			Message: "Nearfield zone detected, diffusion may not be optimal",
		})
	}

	wEnergy := 0.25
	wTime := 0.15
	wReceiver := 0.10
	wMount := 0.10
	wSymmetry := 0.10
	wDistance := 0.10
	wBandMatch := 0.10
	wPathSupport := 0.10

	if r.projectSpaceType == domain.SpaceTypeStudio {
		wSymmetry = 0.20
		nearfieldPenalty *= 2.0
		components.NearfieldPenalty = nearfieldPenalty
	}

	score := wEnergy*energyScore +
		wTime*timeScore +
		wReceiver*receiverScore +
		wMount*mountScore +
		wSymmetry*symmetryScore +
		wDistance*distanceScore +
		wBandMatch*bandMatchScore +
		wPathSupport*pathSupportScore -
		components.ModalPenalty -
		components.NearfieldPenalty

	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	return score, components, reasons
}

func (r *Ranker) normalizeEnergy(totalEnergy float64, count int) float64 {
	if count == 0 {
		return 0
	}
	avgEnergy := totalEnergy / float64(count)
	return math.Min(avgEnergy*100, 1.0)
}

func (r *Ranker) normalizeArrivalTime(arrivalTimeMs float64) float64 {
	if arrivalTimeMs < 10 {
		return 1.0
	}
	if arrivalTimeMs > 50 {
		return 0.3
	}
	return 1.0 - (arrivalTimeMs-10)/40*0.7
}

func (r *Ranker) calculateReceiverScore(candidate *Candidate) float64 {
	if len(candidate.SourceReceiverPairs) == 0 {
		return 0
	}
	return math.Min(float64(len(candidate.SourceReceiverPairs))/3.0, 1.0)
}

func (r *Ranker) calculateSymmetryScore(candidate *Candidate, ctx *CandidateContext) float64 {
	if r.projectSpaceType != domain.SpaceTypeStudio {
		return 0.5
	}

	surfaceKind := candidate.Surface.Kind
	if surfaceKind == domain.SurfaceKindFloor || surfaceKind == domain.SurfaceKindCeiling {
		return 0.5
	}

	return 0.7
}

func (r *Ranker) calculateDistanceScore(candidate *Candidate, ctx *CandidateContext) float64 {
	return 0.6
}

func (r *Ranker) calculateBandMatchScore(candidate *Candidate) float64 {
	if len(candidate.TargetBands) == 0 {
		return 0.5
	}
	return 0.7
}

func (r *Ranker) calculatePathSupportScore(candidate *Candidate) float64 {
	return math.Min(float64(candidate.ReflectionCount)/5.0, 1.0)
}

func (r *Ranker) calculateModalPenalty(candidate *Candidate, ctx *CandidateContext) float64 {
	return 0.1
}

func (r *Ranker) calculateNearfieldPenalty(candidate *Candidate, ctx *CandidateContext) float64 {
	surfaceKind := candidate.Surface.Kind
	if surfaceKind == domain.SurfaceKindFloor || surfaceKind == domain.SurfaceKindCeiling {
		return 0
	}

	if candidate.AvgArrivalTimeMs < 15 {
		return 0.4
	}
	return 0
}
