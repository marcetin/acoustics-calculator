package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TreatmentMode string

const (
	TreatmentModeBalanced        TreatmentMode = "BALANCED"
	TreatmentModeDiffusionFirst TreatmentMode = "DIFFUSION_FIRST"
	TreatmentModeAbsorptionFirst TreatmentMode = "ABSORPTION_FIRST"
)

type BudgetClass string

const (
	BudgetClassLow    BudgetClass = "LOW"
	BudgetClassMedium BudgetClass = "MEDIUM"
	BudgetClassHigh   BudgetClass = "HIGH"
)

type ConstraintSet struct {
	ID                     string         `json:"id"`
	ProjectID              string         `json:"project_id"`
	MaxPanelDepthM         float64        `json:"max_panel_depth_m"`
	SymmetryRequired       bool           `json:"symmetry_required"`
	AllowedSurfaceKindsJSON string        `json:"allowed_surface_kinds_json"`
	ForbiddenSurfaceIDsJSON string        `json:"forbidden_surface_ids_json"`
	PreferredTreatmentMode TreatmentMode  `json:"preferred_treatment_mode"`
	BudgetClass            BudgetClass    `json:"budget_class"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

type CreateConstraintInput struct {
	ProjectID              string         `json:"project_id"`
	MaxPanelDepthM         float64        `json:"max_panel_depth_m"`
	SymmetryRequired       bool           `json:"symmetry_required"`
	AllowedSurfaceKinds    []SurfaceKind  `json:"allowed_surface_kinds"`
	ForbiddenSurfaceIDs    []string       `json:"forbidden_surface_ids"`
	PreferredTreatmentMode TreatmentMode  `json:"preferred_treatment_mode"`
	BudgetClass            BudgetClass    `json:"budget_class"`
}

type UpdateConstraintInput struct {
	MaxPanelDepthM         *float64       `json:"max_panel_depth_m,omitempty"`
	SymmetryRequired       *bool          `json:"symmetry_required,omitempty"`
	AllowedSurfaceKinds    *[]SurfaceKind `json:"allowed_surface_kinds,omitempty"`
	ForbiddenSurfaceIDs    *[]string      `json:"forbidden_surface_ids,omitempty"`
	PreferredTreatmentMode *TreatmentMode `json:"preferred_treatment_mode,omitempty"`
	BudgetClass            *BudgetClass   `json:"budget_class,omitempty"`
}

func NewConstraintSet(input CreateConstraintInput) (*ConstraintSet, error) {
	allowedKindsJSON, err := json.Marshal(input.AllowedSurfaceKinds)
	if err != nil {
		return nil, err
	}
	
	forbiddenIDsJSON, err := json.Marshal(input.ForbiddenSurfaceIDs)
	if err != nil {
		return nil, err
	}
	
	constraint := &ConstraintSet{
		ID:                     uuid.New().String(),
		ProjectID:              input.ProjectID,
		MaxPanelDepthM:         input.MaxPanelDepthM,
		SymmetryRequired:       input.SymmetryRequired,
		AllowedSurfaceKindsJSON: string(allowedKindsJSON),
		ForbiddenSurfaceIDsJSON: string(forbiddenIDsJSON),
		PreferredTreatmentMode: input.PreferredTreatmentMode,
		BudgetClass:            input.BudgetClass,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	
	return constraint, nil
}

func (c *ConstraintSet) Update(input UpdateConstraintInput) error {
	if input.MaxPanelDepthM != nil {
		c.MaxPanelDepthM = *input.MaxPanelDepthM
	}
	if input.SymmetryRequired != nil {
		c.SymmetryRequired = *input.SymmetryRequired
	}
	if input.AllowedSurfaceKinds != nil {
		allowedKindsJSON, err := json.Marshal(input.AllowedSurfaceKinds)
		if err != nil {
			return err
		}
		c.AllowedSurfaceKindsJSON = string(allowedKindsJSON)
	}
	if input.ForbiddenSurfaceIDs != nil {
		forbiddenIDsJSON, err := json.Marshal(input.ForbiddenSurfaceIDs)
		if err != nil {
			return err
		}
		c.ForbiddenSurfaceIDsJSON = string(forbiddenIDsJSON)
	}
	if input.PreferredTreatmentMode != nil {
		c.PreferredTreatmentMode = *input.PreferredTreatmentMode
	}
	if input.BudgetClass != nil {
		c.BudgetClass = *input.BudgetClass
	}
	c.UpdatedAt = time.Now()
	return nil
}

func (c *ConstraintSet) GetAllowedSurfaceKinds() ([]SurfaceKind, error) {
	var kinds []SurfaceKind
	if c.AllowedSurfaceKindsJSON == "" {
		return kinds, nil
	}
	err := json.Unmarshal([]byte(c.AllowedSurfaceKindsJSON), &kinds)
	return kinds, err
}

func (c *ConstraintSet) GetForbiddenSurfaceIDs() ([]string, error) {
	var ids []string
	if c.ForbiddenSurfaceIDsJSON == "" {
		return ids, nil
	}
	err := json.Unmarshal([]byte(c.ForbiddenSurfaceIDsJSON), &ids)
	return ids, err
}

func GetDefaultConstraints(projectID string) CreateConstraintInput {
	return CreateConstraintInput{
		ProjectID:              projectID,
		MaxPanelDepthM:         0.3,
		SymmetryRequired:       false,
		AllowedSurfaceKinds:    []SurfaceKind{SurfaceKindWall, SurfaceKindCeiling, SurfaceKindFloor},
		ForbiddenSurfaceIDs:    []string{},
		PreferredTreatmentMode: TreatmentModeBalanced,
		BudgetClass:            BudgetClassMedium,
	}
}

func (tm TreatmentMode) String() string {
	return string(tm)
}

func (tm TreatmentMode) DisplayName() string {
	switch tm {
	case TreatmentModeBalanced:
		return "Balanced"
	case TreatmentModeDiffusionFirst:
		return "Diffusion First"
	case TreatmentModeAbsorptionFirst:
		return "Absorption First"
	default:
		return string(tm)
	}
}

func ValidTreatmentModes() []TreatmentMode {
	return []TreatmentMode{
		TreatmentModeBalanced,
		TreatmentModeDiffusionFirst,
		TreatmentModeAbsorptionFirst,
	}
}

func (bc BudgetClass) String() string {
	return string(bc)
}

func (bc BudgetClass) DisplayName() string {
	switch bc {
	case BudgetClassLow:
		return "Low"
	case BudgetClassMedium:
		return "Medium"
	case BudgetClassHigh:
		return "High"
	default:
		return string(bc)
	}
}

func ValidBudgetClasses() []BudgetClass {
	return []BudgetClass{
		BudgetClassLow,
		BudgetClassMedium,
		BudgetClassHigh,
	}
}
