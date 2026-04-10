package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AnalysisStatus string

const (
	AnalysisStatusQueued    AnalysisStatus = "queued"
	AnalysisStatusRunning   AnalysisStatus = "running"
	AnalysisStatusCompleted AnalysisStatus = "completed"
	AnalysisStatusFailed    AnalysisStatus = "failed"
)

type AnalysisRun struct {
	ID           string
	ProjectID    string
	Status       AnalysisStatus
	InputsHash   string
	MetricsJSON  string
	StartedAt    *time.Time
	FinishedAt   *time.Time
	ErrorMessage *string
	CreatedAt    time.Time
}

type AnalysisMetrics struct {
	SpeedOfSound     float64             `json:"speed_of_sound"`
	Modes            []Mode             `json:"modes"`
	RTResults        RTResults          `json:"rt_results"`
	Reflections      []ReflectionPath   `json:"reflections"`
	Warnings         []string           `json:"warnings"`
	Summary          AnalysisSummary    `json:"summary"`
}

type Mode struct {
	L     int     `json:"l"`
	M     int     `json:"m"`
	N     int     `json:"n"`
	Freq  float64 `json:"freq_hz"`
	Kind  string  `json:"kind"` // AXIAL, TANGENTIAL, OBLIQUE
}

type RTResults struct {
	Bands []RTBand `json:"bands"`
}

type RTBand struct {
	FrequencyHz             float64 `json:"frequency_hz"`
	SabineT60               float64 `json:"sabine_t60"`
	EyringT60               float64 `json:"eyring_t60"`
	EquivalentAbsorptionArea float64 `json:"equivalent_absorption_area"`
	MeanAbsorption          float64 `json:"mean_absorption"`
}

type ReflectionPath struct {
	SourceID       string   `json:"source_id"`
	ReceiverID     string   `json:"receiver_id"`
	Order          int      `json:"order"`
	SurfaceSequence []string `json:"surface_sequence"`
	PathLengthM    float64  `json:"path_length_m"`
	ArrivalTimeMs  float64  `json:"arrival_time_ms"`
	EnergyEstimate float64  `json:"energy_estimate"`
}

type AnalysisSummary struct {
	RunTimestamp        string `json:"run_timestamp"`
	TotalModes          int    `json:"total_modes"`
	TotalReflections    int    `json:"total_reflections"`
	MaxReflectionOrder  int    `json:"max_reflection_order"`
	WarningCount        int    `json:"warning_count"`
}

type CreateAnalysisRunInput struct {
	ProjectID string
}

func NewAnalysisRun(input CreateAnalysisRunInput) (*AnalysisRun, error) {
	id := uuid.New().String()
	now := time.Now()

	return &AnalysisRun{
		ID:          id,
		ProjectID:   input.ProjectID,
		Status:      AnalysisStatusQueued,
		InputsHash:  "",
		MetricsJSON: "{}",
		StartedAt:   nil,
		FinishedAt:  nil,
		CreatedAt:   now,
	}, nil
}

func (a *AnalysisRun) GetMetrics() (*AnalysisMetrics, error) {
	if a.MetricsJSON == "" {
		return &AnalysisMetrics{}, nil
	}
	var metrics AnalysisMetrics
	if err := json.Unmarshal([]byte(a.MetricsJSON), &metrics); err != nil {
		return nil, err
	}
	return &metrics, nil
}

func (a *AnalysisRun) SetMetrics(metrics *AnalysisMetrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	a.MetricsJSON = string(data)
	return nil
}

func (a *AnalysisRun) MarkRunning() {
	a.Status = AnalysisStatusRunning
	now := time.Now()
	a.StartedAt = &now
}

func (a *AnalysisRun) MarkCompleted() {
	a.Status = AnalysisStatusCompleted
	now := time.Now()
	a.FinishedAt = &now
}

func (a *AnalysisRun) MarkFailed(message string) {
	a.Status = AnalysisStatusFailed
	now := time.Now()
	a.FinishedAt = &now
	a.ErrorMessage = &message
}
