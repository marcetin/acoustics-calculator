package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type ReportService struct {
	projectRepo      *repo.ProjectRepository
	geometryRepo     *repo.GeometryRepository
	surfaceRepo      *repo.SurfaceRepository
	materialRepo     *repo.MaterialRepository
	sourceRepo       *repo.SourceRepository
	receiverRepo     *repo.ReceiverRepository
	constraintRepo   *repo.ConstraintRepository
	analysisRepo     *repo.AnalysisRepository
	placementRepo    *repo.PlacementRepository
	analysisService  *AnalysisService
	placementService *PlacementService
}

func NewReportService(
	projectRepo *repo.ProjectRepository,
	geometryRepo *repo.GeometryRepository,
	surfaceRepo *repo.SurfaceRepository,
	materialRepo *repo.MaterialRepository,
	sourceRepo *repo.SourceRepository,
	receiverRepo *repo.ReceiverRepository,
	constraintRepo *repo.ConstraintRepository,
	analysisRepo *repo.AnalysisRepository,
	placementRepo *repo.PlacementRepository,
	analysisService *AnalysisService,
	placementService *PlacementService,
) *ReportService {
	return &ReportService{
		projectRepo:      projectRepo,
		geometryRepo:     geometryRepo,
		surfaceRepo:      surfaceRepo,
		materialRepo:     materialRepo,
		sourceRepo:       sourceRepo,
		receiverRepo:     receiverRepo,
		constraintRepo:   constraintRepo,
		analysisRepo:     analysisRepo,
		placementRepo:    placementRepo,
		analysisService:  analysisService,
		placementService: placementService,
	}
}

func (s *ReportService) BuildProjectReport(ctx context.Context, projectID string) (*domain.ProjectReport, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	report := &domain.ProjectReport{
		Metadata: s.buildMetadata(project),
	}

	report.Project = &domain.ProjectReportProject{
		ID:              project.ID,
		Name:            project.Name,
		SpaceType:       project.SpaceType,
		GoalType:        project.GoalType,
		Units:           project.Units,
		TemperatureC:    project.TemperatureC,
		HumidityPercent: project.HumidityPercent,
		SpeedOfSound:    project.SpeedOfSound,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
	}

	geometry, err := s.geometryRepo.GetByProjectID(ctx, projectID)
	if err == nil {
		report.Geometry = &domain.ProjectReportGeometry{
			GeometryType: geometry.GeometryType,
			WidthM:       geometry.Width,
			LengthM:      geometry.Length,
			HeightM:      geometry.Height,
			VolumeM3:     geometry.VolumeM3,
		}
	}

	surfaces, err := s.surfaceRepo.ListByProjectID(ctx, projectID)
	if err == nil {
		report.Surfaces = make([]*domain.ProjectReportSurface, len(surfaces))
		for i, surf := range surfaces {
			report.Surfaces[i] = &domain.ProjectReportSurface{
				ID:          surf.ID,
				Name:        surf.Name,
				Kind:        surf.Kind,
				Orientation: "", // Surface doesn't have orientation, compute from normal if needed
				AreaM2:      surf.AreaM2,
				Mountable:   surf.IsMountable,
				MaterialID:  surf.MaterialID,
			}
			if surf.MaterialID != nil {
				mat, err := s.materialRepo.GetByID(ctx, *surf.MaterialID)
				if err == nil {
					report.Surfaces[i].MaterialName = &mat.Name
				}
			}
		}
	}

	materials, err := s.materialRepo.ListByProjectID(ctx, projectID)
	if err == nil {
		report.Materials = make([]*domain.ProjectReportMaterial, len(materials))
		for i, mat := range materials {
			absBands, _ := mat.GetAbsorptionBands()
			absCoeff := make(map[int]float64)
			for _, band := range absBands {
				absCoeff[band.Frequency] = band.Value
			}

			scatBands, _ := mat.GetScatteringBands()
			scatCoeff := make(map[int]float64)
			for _, band := range scatBands {
				scatCoeff[band.Frequency] = band.Value
			}

			report.Materials[i] = &domain.ProjectReportMaterial{
				ID:              mat.ID,
				Name:            mat.Name,
				IsPreset:        mat.IsPreset,
				AbsorptionCoeff: absCoeff,
				ScatteringCoeff: scatCoeff,
			}
		}
	}

	sources, err := s.sourceRepo.ListByProjectID(ctx, projectID)
	if err == nil {
		report.Sources = make([]*domain.ProjectReportSource, len(sources))
		for i, src := range sources {
			report.Sources[i] = &domain.ProjectReportSource{
				ID:        src.ID,
				Name:      src.Name,
				Type:      string(src.Type),
				X:         src.PositionX,
				Y:         src.PositionY,
				Z:         src.PositionZ,
				AimX:      src.AimX,
				AimY:      src.AimY,
				AimZ:      src.AimZ,
				GainDB:    src.GainDB,
				GroupID:   src.GroupID,
				CreatedAt: src.CreatedAt,
				UpdatedAt: src.UpdatedAt,
			}
		}
	}

	receivers, err := s.receiverRepo.ListByProjectID(ctx, projectID)
	if err == nil {
		report.Receivers = make([]*domain.ProjectReportReceiver, len(receivers))
		for i, recv := range receivers {
			report.Receivers[i] = &domain.ProjectReportReceiver{
				ID:        recv.ID,
				Name:      recv.Name,
				Type:      string(recv.Type),
				X:         recv.PositionX,
				Y:         recv.PositionY,
				Z:         recv.PositionZ,
				Weight:    recv.Weight,
				CreatedAt: recv.CreatedAt,
				UpdatedAt: recv.UpdatedAt,
			}
		}
	}

	constraints, err := s.constraintRepo.GetByProjectID(ctx, projectID)
	if err == nil && constraints != nil {
		report.Constraints = &domain.ProjectReportConstraints{
			MaxPanelDepthM:         constraints.MaxPanelDepthM,
			SymmetryRequired:       constraints.SymmetryRequired,
			BudgetClass:            string(constraints.BudgetClass),
			PreferredTreatmentMode: string(constraints.PreferredTreatmentMode),
			AllowedSurfaceKinds:    []domain.SurfaceKind{}, // Constraint doesn't have this field
			ForbiddenSurfaces:      []string{},             // Constraint doesn't have this field
		}
	}

	analysisRun, err := s.analysisRepo.GetLatestByProject(ctx, projectID)
	if err == nil && analysisRun != nil {
		isStale, _ := s.analysisService.IsAnalysisStale(ctx, projectID, analysisRun)

		metrics, _ := analysisRun.GetMetrics()

		report.LatestAnalysis = &domain.ProjectReportAnalysis{
			ID:           analysisRun.ID,
			ProjectID:    analysisRun.ProjectID,
			Status:       string(analysisRun.Status),
			CreatedAt:    analysisRun.CreatedAt,
			InputsHash:   analysisRun.InputsHash,
			SpeedOfSound: metrics.SpeedOfSound,
			Summary:      &metrics.Summary,
			Warnings:     metrics.Warnings,
			Modes:        metrics.Modes,
			RT:           metrics.RTResults,
			Reflections:  metrics.Reflections,
			IsStale:      isStale,
		}
	}

	if analysisRun != nil {
		placements, err := s.placementRepo.ListByAnalysisRun(ctx, analysisRun.ID)
		if err == nil && len(placements) > 0 {
			isStale, _ := s.placementService.ArePlacementsStale(ctx, projectID, placements)

			report.LatestPlacements = &domain.ProjectReportPlacements{
				Candidates: placements,
				IsStale:    isStale,
			}

			report.LatestPlacements.AnalysisRunID = &analysisRun.ID

			summary, _ := s.placementService.SummarizePlacements(ctx, projectID)
			report.LatestPlacements.Summary = summary
		}
	}

	return report, nil
}

func (s *ReportService) BuildPlacementsCSV(ctx context.Context, projectID string) ([]*domain.PlacementCSVRow, error) {
	analysisRun, err := s.analysisRepo.GetLatestByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if analysisRun == nil {
		return []*domain.PlacementCSVRow{}, nil
	}

	placements, err := s.placementRepo.ListByAnalysisRun(ctx, analysisRun.ID)
	if err != nil {
		return nil, err
	}

	if len(placements) == 0 {
		return []*domain.PlacementCSVRow{}, nil
	}

	rows := make([]*domain.PlacementCSVRow, len(placements))

	surfaces, _ := s.surfaceRepo.ListByProjectID(ctx, projectID)
	surfaceMap := make(map[string]*domain.Surface)
	for _, surf := range surfaces {
		surfaceMap[surf.ID] = surf
	}

	for i, p := range placements {
		surf := surfaceMap[p.SurfaceID]
		surfaceName := ""
		if surf != nil {
			surfaceName = surf.Name
		}

		targetBandsArr, _ := p.GetTargetBands()
		targetBandsStrs := make([]string, len(targetBandsArr))
		for j, b := range targetBandsArr {
			targetBandsStrs[j] = fmt.Sprintf("%.0f", b)
		}
		targetBands := strings.Join(targetBandsStrs, ", ")

		reasons, _ := p.GetReasons()
		reasonsArr := make([]string, len(reasons))
		for j, r := range reasons {
			reasonsArr[j] = r.Message
		}
		reasonsSummary := strings.Join(reasonsArr, "; ")

		warnings, _ := p.GetWarnings()
		warningsSummary := strings.Join(warnings, "; ")

		diffuserType := ""
		if p.DiffuserTypeID != nil {
			diffuserType = *p.DiffuserTypeID
		}

		centerX := p.CenterX
		centerY := p.CenterY
		centerZ := p.CenterZ

		panelCount := 1
		coverageArea := 0.0

		bbox, _ := p.GetBoundingBox()
		if bbox != nil {
			coverageArea = (bbox.MaxX - bbox.MinX) * (bbox.MaxY - bbox.MinY)
		}

		rows[i] = &domain.PlacementCSVRow{
			ProjectID:       projectID,
			AnalysisRunID:   analysisRun.ID,
			PlacementID:     p.ID,
			SurfaceID:       p.SurfaceID,
			SurfaceName:     surfaceName,
			Decision:        string(p.Decision),
			Score:           p.Score,
			Confidence:      p.Confidence,
			DiffuserType:    diffuserType,
			Orientation:     p.Orientation,
			CenterX:         centerX,
			CenterY:         centerY,
			CenterZ:         centerZ,
			PanelCount:      panelCount,
			CoverageAreaM2:  coverageArea,
			TargetBands:     targetBands,
			ReasonsSummary:  reasonsSummary,
			WarningsSummary: warningsSummary,
			CreatedAt:       p.CreatedAt,
		}
	}

	return rows, nil
}

func (s *ReportService) GetPrintableReportData(ctx context.Context, projectID string) (*domain.ProjectReport, error) {
	return s.BuildProjectReport(ctx, projectID)
}

func (s *ReportService) buildMetadata(project *domain.Project) *domain.ProjectReportMetadata {
	assumptions := []string{
		"Shoebox geometry only (rectangular rooms)",
		"Image-source reflections limited to 1st and 2nd order",
		"Diffuser layout synthesis is simple grid fill",
		"No post-treatment re-simulation",
		"PUBLIC profile uses simplified heuristics",
	}

	isPublicProfile := project.SpaceType == domain.SpaceTypePublic

	return &domain.ProjectReportMetadata{
		GeneratedAt:       time.Now(),
		ReportVersion:     "1.0",
		SupportedScope:    "SHOEBOX geometry, image-source reflections, diffuser placement recommendations",
		Assumptions:       assumptions,
		IsPublicProfile:   isPublicProfile,
		PublicProvisional: isPublicProfile,
	}
}
