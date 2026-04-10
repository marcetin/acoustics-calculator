package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID                  string    `json:"id"`
	ProjectID           *string   `json:"project_id"`
	Name                string    `json:"name"`
	Category            string    `json:"category"`
	AbsorptionBandsJSON string    `json:"absorption_bands_json"`
	ScatteringBandsJSON string    `json:"scattering_bands_json"`
	IsPreset            bool      `json:"is_preset"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type MaterialBand struct {
	Frequency int     `json:"frequency"`
	Value     float64 `json:"value"`
}

type CreateMaterialInput struct {
	ProjectID       *string        `json:"project_id"`
	Name            string         `json:"name"`
	Category        string         `json:"category"`
	AbsorptionBands []MaterialBand `json:"absorption_bands"`
	ScatteringBands []MaterialBand `json:"scattering_bands"`
	IsPreset        bool           `json:"is_preset"`
}

type UpdateMaterialInput struct {
	Name            *string         `json:"name,omitempty"`
	Category        *string         `json:"category,omitempty"`
	AbsorptionBands *[]MaterialBand `json:"absorption_bands,omitempty"`
	ScatteringBands *[]MaterialBand `json:"scattering_bands,omitempty"`
}

func NewMaterial(input CreateMaterialInput) (*Material, error) {
	absorptionJSON, err := json.Marshal(input.AbsorptionBands)
	if err != nil {
		return nil, err
	}

	scatteringJSON, err := json.Marshal(input.ScatteringBands)
	if err != nil {
		return nil, err
	}

	material := &Material{
		ID:                  uuid.New().String(),
		ProjectID:           input.ProjectID,
		Name:                input.Name,
		Category:            input.Category,
		AbsorptionBandsJSON: string(absorptionJSON),
		ScatteringBandsJSON: string(scatteringJSON),
		IsPreset:            input.IsPreset,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	return material, nil
}

func (m *Material) Update(input UpdateMaterialInput) error {
	if input.Name != nil {
		m.Name = *input.Name
	}
	if input.Category != nil {
		m.Category = *input.Category
	}
	if input.AbsorptionBands != nil {
		absorptionJSON, err := json.Marshal(input.AbsorptionBands)
		if err != nil {
			return err
		}
		m.AbsorptionBandsJSON = string(absorptionJSON)
	}
	if input.ScatteringBands != nil {
		scatteringJSON, err := json.Marshal(input.ScatteringBands)
		if err != nil {
			return err
		}
		m.ScatteringBandsJSON = string(scatteringJSON)
	}
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Material) GetAbsorptionBands() ([]MaterialBand, error) {
	var bands []MaterialBand
	if m.AbsorptionBandsJSON == "" {
		return bands, nil
	}
	err := json.Unmarshal([]byte(m.AbsorptionBandsJSON), &bands)
	return bands, err
}

func (m *Material) GetScatteringBands() ([]MaterialBand, error) {
	var bands []MaterialBand
	if m.ScatteringBandsJSON == "" {
		return bands, nil
	}
	err := json.Unmarshal([]byte(m.ScatteringBandsJSON), &bands)
	return bands, err
}

func (m *Material) IsCustom() bool {
	return !m.IsPreset
}

func GetPresetMaterials() []CreateMaterialInput {
	return []CreateMaterialInput{
		{
			Name:     "Painted Concrete",
			Category: "Wall",
			AbsorptionBands: []MaterialBand{
				{Frequency: 63, Value: 0.01},
				{Frequency: 125, Value: 0.01},
				{Frequency: 250, Value: 0.01},
				{Frequency: 500, Value: 0.02},
				{Frequency: 1000, Value: 0.02},
				{Frequency: 2000, Value: 0.02},
				{Frequency: 4000, Value: 0.03},
			},
			ScatteringBands: []MaterialBand{
				{Frequency: 63, Value: 0.05},
				{Frequency: 125, Value: 0.05},
				{Frequency: 250, Value: 0.05},
				{Frequency: 500, Value: 0.05},
				{Frequency: 1000, Value: 0.05},
				{Frequency: 2000, Value: 0.05},
				{Frequency: 4000, Value: 0.05},
			},
			IsPreset: true,
		},
		{
			Name:     "Gypsum Board",
			Category: "Wall",
			AbsorptionBands: []MaterialBand{
				{Frequency: 63, Value: 0.05},
				{Frequency: 125, Value: 0.05},
				{Frequency: 250, Value: 0.06},
				{Frequency: 500, Value: 0.07},
				{Frequency: 1000, Value: 0.09},
				{Frequency: 2000, Value: 0.10},
				{Frequency: 4000, Value: 0.12},
			},
			ScatteringBands: []MaterialBand{
				{Frequency: 63, Value: 0.05},
				{Frequency: 125, Value: 0.05},
				{Frequency: 250, Value: 0.05},
				{Frequency: 500, Value: 0.05},
				{Frequency: 1000, Value: 0.05},
				{Frequency: 2000, Value: 0.05},
				{Frequency: 4000, Value: 0.05},
			},
			IsPreset: true,
		},
		{
			Name:     "Carpet",
			Category: "Floor",
			AbsorptionBands: []MaterialBand{
				{Frequency: 63, Value: 0.01},
				{Frequency: 125, Value: 0.05},
				{Frequency: 250, Value: 0.10},
				{Frequency: 500, Value: 0.20},
				{Frequency: 1000, Value: 0.35},
				{Frequency: 2000, Value: 0.50},
				{Frequency: 4000, Value: 0.60},
			},
			ScatteringBands: []MaterialBand{
				{Frequency: 63, Value: 0.02},
				{Frequency: 125, Value: 0.02},
				{Frequency: 250, Value: 0.03},
				{Frequency: 500, Value: 0.03},
				{Frequency: 1000, Value: 0.04},
				{Frequency: 2000, Value: 0.04},
				{Frequency: 4000, Value: 0.05},
			},
			IsPreset: true,
		},
		{
			Name:     "Glass",
			Category: "Window",
			AbsorptionBands: []MaterialBand{
				{Frequency: 63, Value: 0.35},
				{Frequency: 125, Value: 0.20},
				{Frequency: 250, Value: 0.12},
				{Frequency: 500, Value: 0.07},
				{Frequency: 1000, Value: 0.04},
				{Frequency: 2000, Value: 0.03},
				{Frequency: 4000, Value: 0.02},
			},
			ScatteringBands: []MaterialBand{
				{Frequency: 63, Value: 0.10},
				{Frequency: 125, Value: 0.10},
				{Frequency: 250, Value: 0.10},
				{Frequency: 500, Value: 0.10},
				{Frequency: 1000, Value: 0.10},
				{Frequency: 2000, Value: 0.10},
				{Frequency: 4000, Value: 0.10},
			},
			IsPreset: true,
		},
		{
			Name:     "Heavy Curtain",
			Category: "Window",
			AbsorptionBands: []MaterialBand{
				{Frequency: 63, Value: 0.05},
				{Frequency: 125, Value: 0.10},
				{Frequency: 250, Value: 0.25},
				{Frequency: 500, Value: 0.55},
				{Frequency: 1000, Value: 0.65},
				{Frequency: 2000, Value: 0.70},
				{Frequency: 4000, Value: 0.70},
			},
			ScatteringBands: []MaterialBand{
				{Frequency: 63, Value: 0.05},
				{Frequency: 125, Value: 0.05},
				{Frequency: 250, Value: 0.05},
				{Frequency: 500, Value: 0.05},
				{Frequency: 1000, Value: 0.05},
				{Frequency: 2000, Value: 0.05},
				{Frequency: 4000, Value: 0.05},
			},
			IsPreset: true,
		},
		{
			Name:     "Wood Panel",
			Category: "Wall",
			AbsorptionBands: []MaterialBand{
				{Frequency: 63, Value: 0.10},
				{Frequency: 125, Value: 0.15},
				{Frequency: 250, Value: 0.20},
				{Frequency: 500, Value: 0.15},
				{Frequency: 1000, Value: 0.10},
				{Frequency: 2000, Value: 0.08},
				{Frequency: 4000, Value: 0.07},
			},
			ScatteringBands: []MaterialBand{
				{Frequency: 63, Value: 0.15},
				{Frequency: 125, Value: 0.15},
				{Frequency: 250, Value: 0.15},
				{Frequency: 500, Value: 0.15},
				{Frequency: 1000, Value: 0.15},
				{Frequency: 2000, Value: 0.15},
				{Frequency: 4000, Value: 0.15},
			},
			IsPreset: true,
		},
	}
}
