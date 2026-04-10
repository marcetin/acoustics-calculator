package service

import (
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestJobService_CreateJob(t *testing.T) {
	job, err := domain.NewJob(domain.CreateJobInput{
		ProjectID: "test-project",
		Kind:      domain.JobKindAnalysis,
	})

	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	if job.ID == "" {
		t.Error("Job ID should not be empty")
	}

	if job.Status != domain.JobStatusQueued {
		t.Errorf("Job status should be QUEUED, got %s", job.Status)
	}

	if job.ProgressPercent != 0 {
		t.Errorf("Job progress should be 0, got %d", job.ProgressPercent)
	}
}

func TestJob_Start(t *testing.T) {
	job, _ := domain.NewJob(domain.CreateJobInput{
		ProjectID: "test-project",
		Kind:      domain.JobKindAnalysis,
	})

	job.Start()

	if job.Status != domain.JobStatusRunning {
		t.Errorf("Job status should be RUNNING, got %s", job.Status)
	}

	if job.StartedAt == nil {
		t.Error("Job StartedAt should not be nil after start")
	}
}

func TestJob_UpdateProgress(t *testing.T) {
	job, _ := domain.NewJob(domain.CreateJobInput{
		ProjectID: "test-project",
		Kind:      domain.JobKindAnalysis,
	})

	job.UpdateProgress("validating_inputs", 10, "Validating inputs")

	if job.Phase != "validating_inputs" {
		t.Errorf("Job phase should be validating_inputs, got %s", job.Phase)
	}

	if job.ProgressPercent != 10 {
		t.Errorf("Job progress should be 10, got %d", job.ProgressPercent)
	}

	if job.Message != "Validating inputs" {
		t.Errorf("Job message should be 'Validating inputs', got %s", job.Message)
	}
}

func TestJob_ProgressMonotonic(t *testing.T) {
	job, _ := domain.NewJob(domain.CreateJobInput{
		ProjectID: "test-project",
		Kind:      domain.JobKindAnalysis,
	})

	job.UpdateProgress("phase1", 10, "Phase 1")
	job.UpdateProgress("phase2", 30, "Phase 2")
	job.UpdateProgress("phase3", 20, "Phase 3")

	if job.ProgressPercent != 30 {
		t.Errorf("Job progress should be 30 (monotonic), got %d", job.ProgressPercent)
	}
}

func TestJob_Complete(t *testing.T) {
	job, _ := domain.NewJob(domain.CreateJobInput{
		ProjectID: "test-project",
		Kind:      domain.JobKindAnalysis,
	})

	job.Start()
	job.Complete()

	if job.Status != domain.JobStatusCompleted {
		t.Errorf("Job status should be COMPLETED, got %s", job.Status)
	}

	if job.ProgressPercent != 100 {
		t.Errorf("Job progress should be 100, got %d", job.ProgressPercent)
	}

	if job.FinishedAt == nil {
		t.Error("Job FinishedAt should not be nil after complete")
	}

	if !job.IsFinished() {
		t.Error("Job should be finished after complete")
	}
}

func TestJob_Fail(t *testing.T) {
	job, _ := domain.NewJob(domain.CreateJobInput{
		ProjectID: "test-project",
		Kind:      domain.JobKindAnalysis,
	})

	job.Start()
	job.Fail("Test error")

	if job.Status != domain.JobStatusFailed {
		t.Errorf("Job status should be FAILED, got %s", job.Status)
	}

	if job.ErrorMessage == nil {
		t.Error("Job ErrorMessage should not be nil after fail")
	}

	if *job.ErrorMessage != "Test error" {
		t.Errorf("Job ErrorMessage should be 'Test error', got %s", *job.ErrorMessage)
	}

	if !job.IsFinished() {
		t.Error("Job should be finished after fail")
	}
}
