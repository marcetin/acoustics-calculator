package domain

import (
	"time"

	"github.com/google/uuid"
)

type GeometryType string

const (
	GeometryTypeShoebox    GeometryType = "SHOEBOX"
	GeometryTypePolyhedral GeometryType = "POLYHEDRAL"
)

type RoomGeometry struct {
	ID           string       `json:"id"`
	ProjectID    string       `json:"project_id"`
	GeometryType GeometryType `json:"geometry_type"`
	Width        float64      `json:"width"`
	Length       float64      `json:"length"`
	Height       float64      `json:"height"`
	VolumeM3     float64      `json:"volume_m3"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type CreateGeometryInput struct {
	ProjectID    string       `json:"project_id"`
	GeometryType GeometryType `json:"geometry_type"`
	Width        float64      `json:"width"`
	Length       float64      `json:"length"`
	Height       float64      `json:"height"`
}

type UpdateGeometryInput struct {
	Width  *float64 `json:"width,omitempty"`
	Length *float64 `json:"length,omitempty"`
	Height *float64 `json:"height,omitempty"`
}

func NewRoomGeometry(input CreateGeometryInput) *RoomGeometry {
	now := time.Now()

	geometry := &RoomGeometry{
		ID:           uuid.New().String(),
		ProjectID:    input.ProjectID,
		GeometryType: input.GeometryType,
		Width:        input.Width,
		Length:       input.Length,
		Height:       input.Height,
		VolumeM3:     calculateVolume(input.Width, input.Length, input.Height),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return geometry
}

func (g *RoomGeometry) Update(input UpdateGeometryInput) {
	if input.Width != nil {
		g.Width = *input.Width
	}
	if input.Length != nil {
		g.Length = *input.Length
	}
	if input.Height != nil {
		g.Height = *input.Height
	}
	g.VolumeM3 = calculateVolume(g.Width, g.Length, g.Height)
	g.UpdatedAt = time.Now()
}

func (g *RoomGeometry) CalculateVolume() float64 {
	return calculateVolume(g.Width, g.Length, g.Height)
}

func calculateVolume(width, length, height float64) float64 {
	return width * length * height
}

func (g *RoomGeometry) IsShoebox() bool {
	return g.GeometryType == GeometryTypeShoebox
}

func (gt GeometryType) String() string {
	return string(gt)
}

func (gt GeometryType) DisplayName() string {
	switch gt {
	case GeometryTypeShoebox:
		return "Rectangular (Shoebox)"
	case GeometryTypePolyhedral:
		return "Polyhedral"
	default:
		return string(gt)
	}
}
