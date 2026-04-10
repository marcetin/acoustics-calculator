package scoring

import (
	"acoustics-calculator/internal/domain"
)

type VetoChecker struct {
	projectSpaceType domain.SpaceType
	constraints      *domain.ConstraintSet
}

func NewVetoChecker(spaceType domain.SpaceType, constraints *domain.ConstraintSet) *VetoChecker {
	return &VetoChecker{
		projectSpaceType: spaceType,
		constraints:      constraints,
	}
}

type VetoResult struct {
	Decision domain.PlacementDecision
	Reasons  []domain.CandidateReason
}

func (v *VetoChecker) Apply(candidate *Candidate, score float64, ctx *CandidateContext) *VetoResult {
	var reasons []domain.CandidateReason

	if !candidate.Surface.IsMountable {
		return &VetoResult{
			Decision: domain.PlacementDecisionReject,
			Reasons: []domain.CandidateReason{
				{Code: "NOT_MOUNTABLE", Message: "Surface is not mountable"},
			},
		}
	}

	patchArea := (candidate.Patch.MaxX - candidate.Patch.MinX) * (candidate.Patch.MaxY - candidate.Patch.MinY)
	if patchArea < 0.25 {
		return &VetoResult{
			Decision: domain.PlacementDecisionReject,
			Reasons: []domain.CandidateReason{
				{Code: "PATCH_TOO_SMALL", Message: "Patch area too small for treatment"},
			},
		}
	}

	if v.constraints != nil && v.constraints.MaxPanelDepthM > 0 {
		_ = v.constraints.MaxPanelDepthM
	}

	if v.isNearfield(candidate) && v.projectSpaceType == domain.SpaceTypeStudio {
		reasons = append(reasons, domain.CandidateReason{
			Code:    "NEARFIELD_STUDIO",
			Message: "Nearfield zone in STUDIO profile - absorber recommended",
		})
		return &VetoResult{
			Decision: domain.PlacementDecisionAbsorberRecommended,
			Reasons:  reasons,
		}
	}

	if v.isFirstReflectionSidewall(candidate) && v.projectSpaceType == domain.SpaceTypeStudio {
		reasons = append(reasons, domain.CandidateReason{
			Code:    "FIRST_REFLECTION_STUDIO",
			Message: "First reflection sidewall in STUDIO - absorber preferred",
		})
		return &VetoResult{
			Decision: domain.PlacementDecisionAbsorberRecommended,
			Reasons:  reasons,
		}
	}

	if score < 0.3 {
		return &VetoResult{
			Decision: domain.PlacementDecisionReject,
			Reasons: []domain.CandidateReason{
				{Code: "LOW_SCORE", Message: "Score too low for viable treatment"},
			},
		}
	}

	if v.projectSpaceType == domain.SpaceTypePublic {
		reasons = append(reasons, domain.CandidateReason{
			Code:    "PUBLIC_PROVISIONAL",
			Message: "PUBLIC profile support is provisional in V1",
		})
	}

	return &VetoResult{
		Decision: domain.PlacementDecisionDiffuser,
		Reasons:  reasons,
	}
}

func (v *VetoChecker) isNearfield(candidate *Candidate) bool {
	return candidate.AvgArrivalTimeMs < 15
}

func (v *VetoChecker) isFirstReflectionSidewall(candidate *Candidate) bool {
	surfaceKind := candidate.Surface.Kind
	if surfaceKind != domain.SurfaceKindWall {
		return false
	}

	return candidate.AvgArrivalTimeMs < 20 && candidate.ReflectionCount > 0
}
