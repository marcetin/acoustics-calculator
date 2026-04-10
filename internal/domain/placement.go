package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PlacementDecision string

const (
	PlacementDecisionDiffuser             PlacementDecision = "DIFFUSER"
	PlacementDecisionAbsorberRecommended PlacementDecision = "ABSORBER_RECOMMENDED"
	PlacementDecisionReject               PlacementDecision = "REJECT"
)

type PlacementCandidate struct {
	ID                   string
	ProjectID            string
	AnalysisRunID        string
	SurfaceID            string
	BoundingBoxJSON      string
	CenterX              float64
	CenterY              float64
	CenterZ              float64
	Orientation          string
	Score                float64
	Confidence           float64
	Decision             PlacementDecision
	DiffuserTypeID       *string
	ReasonsJSON          string
	WarningsJSON         string
	TargetBandsJSON      string
	SourceReceiverPairsJSON string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type BoundingBox struct {
	MinX float64 `json:"min_x"`
	MaxX float64 `json:"max_x"`
	MinY float64 `json:"min_y"`
	MaxY float64 `json:"max_y"`
	MinZ float64 `json:"min_z"`
	MaxZ float64 `json:"max_z"`
}

type CandidateReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CreatePlacementCandidateInput struct {
	ProjectID            string
	AnalysisRunID        string
	SurfaceID            string
	BoundingBox          *BoundingBox
	CenterX              float64
	CenterY              float64
	CenterZ              float64
	Orientation          string
	Score                float64
	Confidence           float64
	Decision             PlacementDecision
	DiffuserTypeID       *string
	Reasons              []CandidateReason
	Warnings             []string
	TargetBands          []float64
	SourceReceiverPairs  []string
}

func NewPlacementCandidate(input CreatePlacementCandidateInput) (*PlacementCandidate, error) {
	id := uuid.New().String()
	now := time.Now()

	boundingBoxJSON, err := json.Marshal(input.BoundingBox)
	if err != nil {
		return nil, err
	}

	reasonsJSON, err := json.Marshal(input.Reasons)
	if err != nil {
		return nil, err
	}

	warningsJSON, err := json.Marshal(input.Warnings)
	if err != nil {
		return nil, err
	}

	targetBandsJSON, err := json.Marshal(input.TargetBands)
	if err != nil {
		return nil, err
	}

	sourceReceiverPairsJSON, err := json.Marshal(input.SourceReceiverPairs)
	if err != nil {
		return nil, err
	}

	return &PlacementCandidate{
		ID:                   id,
		ProjectID:            input.ProjectID,
		AnalysisRunID:        input.AnalysisRunID,
		SurfaceID:            input.SurfaceID,
		BoundingBoxJSON:      string(boundingBoxJSON),
		CenterX:              input.CenterX,
		CenterY:              input.CenterY,
		CenterZ:              input.CenterZ,
		Orientation:          input.Orientation,
		Score:                input.Score,
		Confidence:           input.Confidence,
		Decision:             input.Decision,
		DiffuserTypeID:       input.DiffuserTypeID,
		ReasonsJSON:          string(reasonsJSON),
		WarningsJSON:         string(warningsJSON),
		TargetBandsJSON:      string(targetBandsJSON),
		SourceReceiverPairsJSON: string(sourceReceiverPairsJSON),
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func (p *PlacementCandidate) GetBoundingBox() (*BoundingBox, error) {
	if p.BoundingBoxJSON == "" {
		return nil, nil
	}
	var box BoundingBox
	if err := json.Unmarshal([]byte(p.BoundingBoxJSON), &box); err != nil {
		return nil, err
	}
	return &box, nil
}

func (p *PlacementCandidate) GetReasons() ([]CandidateReason, error) {
	if p.ReasonsJSON == "" {
		return []CandidateReason{}, nil
	}
	var reasons []CandidateReason
	if err := json.Unmarshal([]byte(p.ReasonsJSON), &reasons); err != nil {
		return nil, err
	}
	return reasons, nil
}

func (p *PlacementCandidate) GetWarnings() ([]string, error) {
	if p.WarningsJSON == "" {
		return []string{}, nil
	}
	var warnings []string
	if err := json.Unmarshal([]byte(p.WarningsJSON), &warnings); err != nil {
		return nil, err
	}
	return warnings, nil
}

func (p *PlacementCandidate) GetTargetBands() ([]float64, error) {
	if p.TargetBandsJSON == "" {
		return []float64{}, nil
	}
	var bands []float64
	if err := json.Unmarshal([]byte(p.TargetBandsJSON), &bands); err != nil {
		return nil, err
	}
	return bands, nil
}

func (p *PlacementCandidate) GetSourceReceiverPairs() ([]string, error) {
	if p.SourceReceiverPairsJSON == "" {
		return []string{}, nil
	}
	var pairs []string
	if err := json.Unmarshal([]byte(p.SourceReceiverPairsJSON), &pairs); err != nil {
		return nil, err
	}
	return pairs, nil
}

type PlacementSummary struct {
	CandidateCount            int
	DiffuserCount            int
	AbsorberRecommendedCount  int
	RejectCount              int
	TopSurfaceNames          []string
	WarningCount             int
}

type LayoutResult struct {
	BoundingBox    *BoundingBox
	PanelCenters   []Vec3
	PanelCount     int
	Orientation    string
	CoverageAreaM2 float64
}

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
