package domain

import (
	"time"

	"github.com/google/uuid"
)

type SpaceType string

const (
	SpaceTypeStudio SpaceType = "STUDIO"
	SpaceTypeHifi   SpaceType = "HIFI"
	SpaceTypePublic SpaceType = "PUBLIC"
)

type GoalType string

const (
	GoalTypeMixing     GoalType = "MIXING"
	GoalTypeMastering  GoalType = "MASTERING"
	GoalTypeHifiStereo GoalType = "HIFI_STEREO"
	GoalTypeSpeech     GoalType = "SPEECH"
	GoalTypeMusic      GoalType = "MUSIC"
	GoalTypeMultipurpose GoalType = "MULTIPURPOSE"
)

type Units string

const (
	UnitsMetric Units = "METERS"
	UnitsImperial Units = "FEET"
)

type Project struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	SpaceType       SpaceType `json:"space_type"`
	GoalType        GoalType  `json:"goal_type"`
	Units           Units     `json:"units"`
	TemperatureC    float64   `json:"temperature_c"`
	HumidityPercent float64   `json:"humidity_percent"`
	SpeedOfSound    float64   `json:"speed_of_sound"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateProjectInput struct {
	Name            string    `json:"name"`
	SpaceType       SpaceType `json:"space_type"`
	GoalType        GoalType  `json:"goal_type"`
	Units           Units     `json:"units"`
	TemperatureC    float64   `json:"temperature_c"`
	HumidityPercent float64   `json:"humidity_percent"`
	SpeedOfSound    float64   `json:"speed_of_sound"`
}

type UpdateProjectInput struct {
	Name            *string    `json:"name,omitempty"`
	SpaceType       *SpaceType `json:"space_type,omitempty"`
	GoalType        *GoalType  `json:"goal_type,omitempty"`
	Units           *Units     `json:"units,omitempty"`
	TemperatureC    *float64   `json:"temperature_c,omitempty"`
	HumidityPercent *float64   `json:"humidity_percent,omitempty"`
	SpeedOfSound    *float64   `json:"speed_of_sound,omitempty"`
}

func NewProject(input CreateProjectInput) *Project {
	now := time.Now()
	
	return &Project{
		ID:              uuid.New().String(),
		Name:            input.Name,
		SpaceType:       input.SpaceType,
		GoalType:        input.GoalType,
		Units:           input.Units,
		TemperatureC:    input.TemperatureC,
		HumidityPercent: input.HumidityPercent,
		SpeedOfSound:    input.SpeedOfSound,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (p *Project) Update(input UpdateProjectInput) {
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.SpaceType != nil {
		p.SpaceType = *input.SpaceType
	}
	if input.GoalType != nil {
		p.GoalType = *input.GoalType
	}
	if input.Units != nil {
		p.Units = *input.Units
	}
	if input.TemperatureC != nil {
		p.TemperatureC = *input.TemperatureC
	}
	if input.HumidityPercent != nil {
		p.HumidityPercent = *input.HumidityPercent
	}
	if input.SpeedOfSound != nil {
		p.SpeedOfSound = *input.SpeedOfSound
	}
	p.UpdatedAt = time.Now()
}

func (s SpaceType) String() string {
	return string(s)
}

func (s SpaceType) DisplayName() string {
	switch s {
	case SpaceTypeStudio:
		return "Studio / Control Room"
	case SpaceTypeHifi:
		return "Hi-fi Listening Room"
	case SpaceTypePublic:
		return "Public Space"
	default:
		return string(s)
	}
}

func (s SpaceType) Description() string {
	switch s {
	case SpaceTypeStudio:
		return "Critical listening environment with symmetry and early reflection control"
	case SpaceTypeHifi:
		return "Optimized listening room for imaging and spaciousness"
	case SpaceTypePublic:
		return "Audience-focused space with coverage and clarity requirements"
	default:
		return ""
	}
}

func (g GoalType) String() string {
	return string(g)
}

func (g GoalType) DisplayName() string {
	switch g {
	case GoalTypeMixing:
		return "Mixing"
	case GoalTypeMastering:
		return "Mastering"
	case GoalTypeHifiStereo:
		return "Hi-fi Stereo"
	case GoalTypeSpeech:
		return "Speech"
	case GoalTypeMusic:
		return "Music"
	case GoalTypeMultipurpose:
		return "Multipurpose"
	default:
		return string(g)
	}
}

func (g GoalType) Description() string {
	switch g {
	case GoalTypeMixing:
		return "Music mixing and production"
	case GoalTypeMastering:
		return "Final audio mastering"
	case GoalTypeHifiStereo:
		return "Critical stereo listening"
	case GoalTypeSpeech:
		return "Speech intelligibility and clarity"
	case GoalTypeMusic:
		return "Music performance and enjoyment"
	case GoalTypeMultipurpose:
		return "Multi-use space flexibility"
	default:
		return ""
	}
}

func ValidSpaceTypes() []SpaceType {
	return []SpaceType{SpaceTypeStudio, SpaceTypeHifi, SpaceTypePublic}
}

func ValidGoalTypes() []GoalType {
	return []GoalType{
		GoalTypeMixing,
		GoalTypeMastering,
		GoalTypeHifiStereo,
		GoalTypeSpeech,
		GoalTypeMusic,
		GoalTypeMultipurpose,
	}
}

func ValidUnits() []Units {
	return []Units{UnitsMetric, UnitsImperial}
}
