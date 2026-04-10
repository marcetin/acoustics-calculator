package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobKind string

const (
	JobKindAnalysis  JobKind = "ANALYSIS"
	JobKindPlacement JobKind = "PLACEMENTS"
)

type JobStatus string

const (
	JobStatusQueued    JobStatus = "QUEUED"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

type Job struct {
	ID              string
	ProjectID       string
	Kind            JobKind
	Status          JobStatus
	Phase           string
	ProgressPercent int
	Message         string
	ErrorMessage    *string
	CreatedAt       time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
}

type CreateJobInput struct {
	ProjectID string
	Kind      JobKind
}

func NewJob(input CreateJobInput) (*Job, error) {
	id := uuid.New().String()
	now := time.Now()

	return &Job{
		ID:              id,
		ProjectID:       input.ProjectID,
		Kind:            input.Kind,
		Status:          JobStatusQueued,
		Phase:           "queued",
		ProgressPercent: 0,
		Message:         "Job queued",
		CreatedAt:       now,
	}, nil
}

func (j *Job) Start() {
	now := time.Now()
	j.Status = JobStatusRunning
	j.Phase = "starting"
	j.StartedAt = &now
	j.Message = "Job started"
}

func (j *Job) UpdateProgress(phase string, progressPercent int, message string) {
	j.Phase = phase
	if progressPercent > j.ProgressPercent && progressPercent <= 100 {
		j.ProgressPercent = progressPercent
	}
	j.Message = message
}

func (j *Job) Complete() {
	now := time.Now()
	j.Status = JobStatusCompleted
	j.Phase = "completed"
	j.ProgressPercent = 100
	j.Message = "Job completed successfully"
	j.FinishedAt = &now
}

func (j *Job) Fail(errorMessage string) {
	now := time.Now()
	j.Status = JobStatusFailed
	j.Phase = "failed"
	j.Message = "Job failed"
	j.ErrorMessage = &errorMessage
	j.FinishedAt = &now
}

func (j *Job) IsFinished() bool {
	return j.Status == JobStatusCompleted || j.Status == JobStatusFailed
}
