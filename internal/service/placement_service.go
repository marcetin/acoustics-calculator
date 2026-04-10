package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/acoustics/diffuser"
	"acoustics-calculator/internal/acoustics/scoring"
	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type PlacementService struct {
	placementRepo  *repo.PlacementRepository
	projectRepo    *repo.ProjectRepository
	geometryRepo   *repo.GeometryRepository
	surfaceRepo    *repo.SurfaceRepository
	analysisRepo   *repo.AnalysisRepository
	diffuserRepo   *repo.DiffuserRepository
	constraintRepo *repo.ConstraintRepository
}

func NewPlacementService(
	placementRepo *repo.PlacementRepository,
	projectRepo *repo.ProjectRepository,
	geometryRepo *repo.GeometryRepository,
	surfaceRepo *repo.SurfaceRepository,
	analysisRepo *repo.AnalysisRepository,
	diffuserRepo *repo.DiffuserRepository,
	constraintRepo *repo.ConstraintRepository,
) *PlacementService {
	return &PlacementService{
		placementRepo:  placementRepo,
		projectRepo:    projectRepo,
		geometryRepo:   geometryRepo,
		surfaceRepo:    surfaceRepo,
		analysisRepo:   analysisRepo,
		diffuserRepo:   diffuserRepo,
		constraintRepo: constraintRepo,
	}
}

type PlacementReadiness struct {
	Ready    bool
	Reason   string
	Warnings []string
}

func (s *PlacementService) CanGeneratePlacements(ctx context.Context, projectID string) (*PlacementReadiness, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return &PlacementReadiness{Ready: false, Reason: "Project not found"}, nil
	}

	geometry, err := s.geometryRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return &PlacementReadiness{Ready: false, Reason: "Geometry not found"}, nil
	}

	if geometry.GeometryType != domain.GeometryTypeShoebox {
		return &PlacementReadiness{
			Ready:  false,
			Reason: "PHASE 4 supports SHOEBOX geometry only",
		}, nil
	}

	latestRun, err := s.analysisRepo.GetLatestByProject(ctx, projectID)
	if err != nil || latestRun == nil {
		return &PlacementReadiness{Ready: false, Reason: "No completed analysis run found"}, nil
	}

	if latestRun.Status != domain.AnalysisStatusCompleted {
		return &PlacementReadiness{Ready: false, Reason: "Latest analysis run not completed"}, nil
	}

	metrics, err := latestRun.GetMetrics()
	if err != nil {
		return &PlacementReadiness{Ready: false, Reason: "Failed to parse analysis metrics"}, nil
	}

	if len(metrics.Reflections) == 0 {
		return &PlacementReadiness{Ready: false, Reason: "No reflection data in analysis"}, nil
	}

	warnings := []string{}
	if project.SpaceType == domain.SpaceTypePublic {
		warnings = append(warnings, "PUBLIC profile support is provisional in V1")
	}

	return &PlacementReadiness{Ready: true, Warnings: warnings}, nil
}

func (s *PlacementService) GeneratePlacements(ctx context.Context, projectID string) ([]*domain.PlacementCandidate, error) {
	readiness, err := s.CanGeneratePlacements(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !readiness.Ready {
		return nil, fmt.Errorf("prerequisites not met: %s", readiness.Reason)
	}

	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	geometry, err := s.geometryRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	surfaces, err := s.surfaceRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	surfaceSlice := make([]domain.Surface, len(surfaces))
	for i, s := range surfaces {
		surfaceSlice[i] = *s
	}

	latestRun, err := s.analysisRepo.GetLatestByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	metrics, err := latestRun.GetMetrics()
	if err != nil {
		return nil, err
	}

	constraints, err := s.constraintRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		constraints = nil
	}

	diffusers, err := s.diffuserRepo.ListPresets(ctx)
	if err != nil {
		return nil, err
	}

	diffuserPointers := make([]*domain.DiffuserType, len(diffusers))
	for i := range diffusers {
		diffuserPointers[i] = &diffusers[i]
	}

	catalog := diffuser.NewCatalog(diffuserPointers)

	ctxData := &scoring.CandidateContext{
		Project:         project,
		Geometry:        geometry,
		Constraints:     constraints,
		LatestAnalysis:  metrics,
		Surfaces:        surfaceSlice,
		DiffuserCatalog: catalog,
	}

	extractor := scoring.NewCandidateExtractor()
	candidates := extractor.Extract(ctxData)

	if err := s.placementRepo.DeleteByAnalysisRun(ctx, latestRun.ID); err != nil {
		return nil, err
	}

	ranker := scoring.NewRanker(project.SpaceType)
	vetoChecker := scoring.NewVetoChecker(project.SpaceType, constraints)
	layoutSynthesizer := diffuser.NewLayoutSynthesizer()

	var placementCandidates []*domain.PlacementCandidate

	for _, candidate := range candidates {
		score, _, reasons := ranker.Score(&candidate, ctxData)

		vetoResult := vetoChecker.Apply(&candidate, score, ctxData)
		decision := vetoResult.Decision
		allReasons := append(reasons, vetoResult.Reasons...)

		var diffuserTypeID *string
		var warnings []string

		if decision == domain.PlacementDecisionDiffuser {
			patchArea := (candidate.Patch.MaxX - candidate.Patch.MinX) * (candidate.Patch.MaxY - candidate.Patch.MinY)
			maxDepth := 0.3
			if constraints != nil && constraints.MaxPanelDepthM > 0 {
				maxDepth = constraints.MaxPanelDepthM
			}

			compatibleDiffuser := catalog.FindCompatible(candidate.TargetBands, maxDepth, patchArea, candidate.Orientation)
			if compatibleDiffuser != nil {
				diffuserTypeID = &compatibleDiffuser.ID
				_ = layoutSynthesizer.SynthesizeLayout(candidate.Patch, compatibleDiffuser, candidate.Orientation)
			} else {
				decision = domain.PlacementDecisionAbsorberRecommended
				allReasons = append(allReasons, domain.CandidateReason{
					Code:    "NO_COMPATIBLE_DIFFUSER",
					Message: "No compatible diffuser found for constraints",
				})
			}
		}

		if project.SpaceType == domain.SpaceTypePublic {
			warnings = append(warnings, "PUBLIC profile heuristics only in this phase")
		}

		confidence := score
		if decision != domain.PlacementDecisionDiffuser {
			confidence = score * 0.7
		}

		input := domain.CreatePlacementCandidateInput{
			ProjectID:           projectID,
			AnalysisRunID:       latestRun.ID,
			SurfaceID:           candidate.SurfaceID,
			BoundingBox:         candidate.Patch,
			CenterX:             candidate.Center.X,
			CenterY:             candidate.Center.Y,
			CenterZ:             candidate.Center.Z,
			Orientation:         candidate.Orientation,
			Score:               score,
			Confidence:          confidence,
			Decision:            decision,
			DiffuserTypeID:      diffuserTypeID,
			Reasons:             allReasons,
			Warnings:            warnings,
			TargetBands:         candidate.TargetBands,
			SourceReceiverPairs: candidate.SourceReceiverPairs,
		}

		placementCandidate, err := domain.NewPlacementCandidate(input)
		if err != nil {
			return nil, err
		}

		if err := s.placementRepo.Create(ctx, placementCandidate); err != nil {
			return nil, err
		}

		placementCandidates = append(placementCandidates, placementCandidate)
	}

	return placementCandidates, nil
}

func (s *PlacementService) GetLatestPlacements(ctx context.Context, projectID string) ([]*domain.PlacementCandidate, error) {
	latestRun, err := s.analysisRepo.GetLatestByProject(ctx, projectID)
	if err != nil || latestRun == nil {
		return []*domain.PlacementCandidate{}, nil
	}

	return s.placementRepo.ListByAnalysisRun(ctx, latestRun.ID)
}

func (s *PlacementService) SummarizePlacements(ctx context.Context, projectID string) (*domain.PlacementSummary, error) {
	candidates, err := s.GetLatestPlacements(ctx, projectID)
	if err != nil {
		return nil, err
	}

	summary := &domain.PlacementSummary{
		CandidateCount: len(candidates),
	}

	surfaceCount := make(map[string]int)
	totalWarnings := 0

	for _, c := range candidates {
		switch c.Decision {
		case domain.PlacementDecisionDiffuser:
			summary.DiffuserCount++
		case domain.PlacementDecisionAbsorberRecommended:
			summary.AbsorberRecommendedCount++
		case domain.PlacementDecisionReject:
			summary.RejectCount++
		}

		surfaceCount[c.SurfaceID]++

		warnings, err := c.GetWarnings()
		if err == nil {
			totalWarnings += len(warnings)
		}
	}

	summary.WarningCount = totalWarnings

	var topSurfaces []string
	for surfaceID, count := range surfaceCount {
		if count >= 1 {
			topSurfaces = append(topSurfaces, surfaceID)
		}
	}
	if len(topSurfaces) > 5 {
		topSurfaces = topSurfaces[:5]
	}
	summary.TopSurfaceNames = topSurfaces

	return summary, nil
}
