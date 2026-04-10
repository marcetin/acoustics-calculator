package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type JobService struct {
	jobRepo *repo.JobRepository
}

func NewJobService(jobRepo *repo.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

func (s *JobService) CreateJob(ctx context.Context, projectID string, kind domain.JobKind) (*domain.Job, error) {
	input := domain.CreateJobInput{
		ProjectID: projectID,
		Kind:      kind,
	}
	job, err := domain.NewJob(input)
	if err != nil {
		return nil, err
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *JobService) StartJob(ctx context.Context, job *domain.Job) error {
	job.Start()
	return s.jobRepo.Update(ctx, job)
}

func (s *JobService) UpdateJobProgress(ctx context.Context, jobID string, phase string, progressPercent int, message string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	job.UpdateProgress(phase, progressPercent, message)
	return s.jobRepo.Update(ctx, job)
}

func (s *JobService) CompleteJob(ctx context.Context, jobID string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	job.Complete()
	return s.jobRepo.Update(ctx, job)
}

func (s *JobService) FailJob(ctx context.Context, jobID string, errorMessage string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	job.Fail(errorMessage)
	return s.jobRepo.Update(ctx, job)
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (*domain.Job, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

func (s *JobService) GetLatestByProject(ctx context.Context, projectID string) (*domain.Job, error) {
	return s.jobRepo.GetLatestByProject(ctx, projectID)
}

func (s *JobService) GetLatestActiveByProject(ctx context.Context, projectID string) (*domain.Job, error) {
	return s.jobRepo.GetLatestActiveByProject(ctx, projectID)
}

func (s *JobService) GetLatestByProjectAndKind(ctx context.Context, projectID string, kind domain.JobKind) (*domain.Job, error) {
	return s.jobRepo.GetLatestByProjectAndKind(ctx, projectID, kind)
}

func (s *JobService) RunAnalysisJob(ctx context.Context, projectID string, analysisService *AnalysisService) (string, error) {
	job, err := s.CreateJob(ctx, projectID, domain.JobKindAnalysis)
	if err != nil {
		return "", err
	}

	if err := s.StartJob(ctx, job); err != nil {
		return "", err
	}

	go func() {
		if err := s.executeAnalysisJob(ctx, job.ID, projectID, analysisService); err != nil {
			s.FailJob(ctx, job.ID, err.Error())
		}
	}()

	return job.ID, nil
}

func (s *JobService) executeAnalysisJob(ctx context.Context, jobID string, projectID string, analysisService *AnalysisService) error {
	s.UpdateJobProgress(ctx, jobID, "validating_inputs", 10, "Validating inputs")

	_, err := analysisService.RunAnalysis(ctx, projectID)
	if err != nil {
		return err
	}

	s.UpdateJobProgress(ctx, jobID, "persisting_results", 90, "Persisting results")

	s.CompleteJob(ctx, jobID)

	return nil
}

func (s *JobService) RunPlacementsJob(ctx context.Context, projectID string, placementService *PlacementService) (string, error) {
	job, err := s.CreateJob(ctx, projectID, domain.JobKindPlacement)
	if err != nil {
		return "", err
	}

	if err := s.StartJob(ctx, job); err != nil {
		return "", err
	}

	go func() {
		if err := s.executePlacementsJob(ctx, job.ID, projectID, placementService); err != nil {
			s.FailJob(ctx, job.ID, err.Error())
		}
	}()

	return job.ID, nil
}

func (s *JobService) executePlacementsJob(ctx context.Context, jobID string, projectID string, placementService *PlacementService) error {
	s.UpdateJobProgress(ctx, jobID, "validating_inputs", 10, "Validating prerequisites")

	s.UpdateJobProgress(ctx, jobID, "loading_analysis_context", 20, "Loading analysis context")

	s.UpdateJobProgress(ctx, jobID, "extracting_candidates", 30, "Extracting candidates")

	s.UpdateJobProgress(ctx, jobID, "scoring_candidates", 50, "Scoring candidates")

	s.UpdateJobProgress(ctx, jobID, "applying_veto", 60, "Applying veto rules")

	s.UpdateJobProgress(ctx, jobID, "assigning_diffusers", 70, "Assigning diffuser types")

	s.UpdateJobProgress(ctx, jobID, "synthesizing_layouts", 80, "Synthesizing layouts")

	s.UpdateJobProgress(ctx, jobID, "persisting_results", 90, "Persisting results")

	_, err := placementService.GeneratePlacements(ctx, projectID)
	if err != nil {
		return fmt.Errorf("placement generation failed: %w", err)
	}

	s.CompleteJob(ctx, jobID)

	return nil
}
