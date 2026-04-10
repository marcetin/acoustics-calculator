package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DiffuserCategory string

const (
	DiffuserCategoryQRD1D   DiffuserCategory = "QRD_1D"
	DiffuserCategoryQRD2D   DiffuserCategory = "QRD_2D"
	DiffuserCategoryCustom DiffuserCategory = "CUSTOM"
)

type ScatterAxis string

const (
	ScatterAxisHorizontal ScatterAxis = "HORIZONTAL"
	ScatterAxisVertical   ScatterAxis = "VERTICAL"
	ScatterAxisBoth       ScatterAxis = "BOTH"
)

type DiffuserType struct {
	ID                string
	Name              string
	Category          DiffuserCategory
	PanelWidthM       float64
	PanelHeightM      float64
	PanelDepthM       float64
	MinEffectiveHz    float64
	MaxEffectiveHz    float64
	ScatterAxis       ScatterAxis
	MountOrientationJSON *string
	CostClass         *string
	IsPreset          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CreateDiffuserTypeInput struct {
	Name              string
	Category          DiffuserCategory
	PanelWidthM       float64
	PanelHeightM      float64
	PanelDepthM       float64
	MinEffectiveHz    float64
	MaxEffectiveHz    float64
	ScatterAxis       ScatterAxis
	MountOrientation  *string
	CostClass         *string
	IsPreset          bool
}

func NewDiffuserType(input CreateDiffuserTypeInput) (*DiffuserType, error) {
	id := uuid.New().String()
	now := time.Now()

	var mountOrientationJSON *string
	if input.MountOrientation != nil {
		data, err := json.Marshal(input.MountOrientation)
		if err != nil {
			return nil, err
		}
		str := string(data)
		mountOrientationJSON = &str
	}

	return &DiffuserType{
		ID:                id,
		Name:              input.Name,
		Category:          input.Category,
		PanelWidthM:       input.PanelWidthM,
		PanelHeightM:      input.PanelHeightM,
		PanelDepthM:       input.PanelDepthM,
		MinEffectiveHz:    input.MinEffectiveHz,
		MaxEffectiveHz:    input.MaxEffectiveHz,
		ScatterAxis:       input.ScatterAxis,
		MountOrientationJSON: mountOrientationJSON,
		CostClass:         input.CostClass,
		IsPreset:          input.IsPreset,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}

func (d *DiffuserType) GetMountOrientation() (map[string]interface{}, error) {
	if d.MountOrientationJSON == nil {
		return nil, nil
	}
	var orientation map[string]interface{}
	if err := json.Unmarshal([]byte(*d.MountOrientationJSON), &orientation); err != nil {
		return nil, err
	}
	return orientation, nil
}
