package service

import (
	"context"
	"fmt"
	"time"

	"acoustics-calculator/internal/acoustics"
	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type AnalysisService struct {
	analysisRepo *repo.AnalysisRepository
	projectRepo  *repo.ProjectRepository
	geometryRepo *repo.GeometryRepository
	surfaceRepo  *repo.SurfaceRepository
	sourceRepo   *repo.SourceRepository
	receiverRepo *repo.ReceiverRepository
	materialRepo *repo.MaterialRepository
}

func NewAnalysisService(
	analysisRepo *repo.AnalysisRepository,
	projectRepo *repo.ProjectRepository,
	geometryRepo *repo.GeometryRepository,
	surfaceRepo *repo.SurfaceRepository,
	sourceRepo *repo.SourceRepository,
	receiverRepo *repo.ReceiverRepository,
	materialRepo *repo.MaterialRepository,
) *AnalysisService {
	return &AnalysisService{
		analysisRepo: analysisRepo,
		projectRepo:  projectRepo,
		geometryRepo: geometryRepo,
		surfaceRepo:  surfaceRepo,
		sourceRepo:   sourceRepo,
		receiverRepo: receiverRepo,
		materialRepo: materialRepo,
	}
}

func (s *AnalysisService) RunAnalysis(ctx context.Context, projectID string) (*domain.AnalysisRun, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	geometry, err := s.geometryRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("geometry not found: %w", err)
	}

	if geometry.GeometryType != domain.GeometryTypeShoebox {
		return nil, fmt.Errorf("unsupported geometry type: %s", geometry.GeometryType)
	}

	surfaces, err := s.surfaceRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load surfaces: %w", err)
	}

	if len(surfaces) == 0 {
		return nil, fmt.Errorf("no surfaces found")
	}

	sources, err := s.sourceRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load sources: %w", err)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no sources found")
	}

	receivers, err := s.receiverRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load receivers: %w", err)
	}

	if len(receivers) == 0 {
		return nil, fmt.Errorf("no receivers found")
	}

	input := domain.CreateAnalysisRunInput{ProjectID: projectID}
	analysisRun, err := domain.NewAnalysisRun(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create analysis run: %w", err)
	}

	if err := s.analysisRepo.Create(ctx, analysisRun); err != nil {
		return nil, fmt.Errorf("failed to save analysis run: %w", err)
	}

	analysisRun.MarkRunning()
	if err := s.analysisRepo.Update(ctx, analysisRun); err != nil {
		return nil, fmt.Errorf("failed to update analysis run: %w", err)
	}

	surfaceSlice := make([]domain.Surface, len(surfaces))
	for i, s := range surfaces {
		surfaceSlice[i] = *s
	}

	sourceSlice := make([]domain.Source, len(sources))
	for i, s := range sources {
		sourceSlice[i] = *s
	}

	receiverSlice := make([]domain.Receiver, len(receivers))
	for i, r := range receivers {
		receiverSlice[i] = *r
	}

	metrics, err := s.calculateMetrics(ctx, project, geometry, surfaceSlice, sourceSlice, receiverSlice)
	if err != nil {
		analysisRun.MarkFailed(err.Error())
		s.analysisRepo.Update(ctx, analysisRun)
		return analysisRun, err
	}

	if err := analysisRun.SetMetrics(metrics); err != nil {
		analysisRun.MarkFailed("failed to serialize metrics")
		s.analysisRepo.Update(ctx, analysisRun)
		return analysisRun, fmt.Errorf("failed to serialize metrics: %w", err)
	}

	analysisRun.MarkCompleted()
	if err := s.analysisRepo.Update(ctx, analysisRun); err != nil {
		return nil, fmt.Errorf("failed to update analysis run: %w", err)
	}

	return analysisRun, nil
}

func (s *AnalysisService) calculateMetrics(ctx context.Context, project *domain.Project, geometry *domain.RoomGeometry, surfaces []domain.Surface, sources []domain.Source, receivers []domain.Receiver) (*domain.AnalysisMetrics, error) {
	var warnings []string

	c := acoustics.CalculateSpeedOfSound(project.TemperatureC)
	warnings = append(warnings, fmt.Sprintf("Speed of sound calculated as %.1f m/s (humidity ignored)", c))

	volume := geometry.Width * geometry.Length * geometry.Height

	maxFreq := 300.0
	modalSolver := acoustics.NewModalSolver(geometry.Width, geometry.Length, geometry.Height, c)
	modes := modalSolver.CalculateModes(maxFreq)

	rtEstimator := acoustics.NewRTEstimator(volume, c)

	materials, err := s.materialRepo.ListPresets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load materials: %w", err)
	}

	materialMap := make(map[string]*domain.Material)
	for _, material := range materials {
		materialMap[material.ID] = material
	}

	rtResults, rtWarnings := rtEstimator.CalculateRT(surfaces, materialMap)
	warnings = append(warnings, rtWarnings...)

	imageSourceSolver := acoustics.NewImageSourceSolver(geometry.Width, geometry.Length, geometry.Height, c)
	reflections, refWarnings := imageSourceSolver.CalculateReflections(sources, receivers, surfaces)
	warnings = append(warnings, refWarnings...)

	summary := domain.AnalysisSummary{
		RunTimestamp:       time.Now().Format(time.RFC3339),
		TotalModes:         len(modes),
		TotalReflections:   len(reflections),
		MaxReflectionOrder: 2,
		WarningCount:       len(warnings),
	}

	return &domain.AnalysisMetrics{
		SpeedOfSound: c,
		Modes:        modes,
		RTResults:    rtResults,
		Reflections:  reflections,
		Warnings:     warnings,
		Summary:      summary,
	}, nil
}

func (s *AnalysisService) GetLatestRun(ctx context.Context, projectID string) (*domain.AnalysisRun, error) {
	return s.analysisRepo.GetLatestByProject(ctx, projectID)
}
