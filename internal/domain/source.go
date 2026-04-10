package domain

import (
	"time"

	"github.com/google/uuid"
)

type SourceType string

const (
	SourceTypeSpeaker    SourceType = "SPEAKER"
	SourceTypePointSource SourceType = "POINT_SOURCE"
	SourceTypeLineArray  SourceType = "LINE_ARRAY"
)

type Source struct {
        ID        string     `json:"id"`
        ProjectID string     `json:"project_id"`
        Name      string     `json:"name"`
        Type      SourceType `json:"type"`
        PositionX float64    `json:"position_x"`
        PositionY float64    `json:"position_y"`
        PositionZ float64    `json:"position_z"`
        AimX      float64    `json:"aim_x"`
        AimY      float64    `json:"aim_y"`
        AimZ      float64    `json:"aim_z"`
        GainDB    float64    `json:"gain_db"`
        GroupID   *string    `json:"group_id"`
        CreatedAt time.Time  `json:"created_at"`
        UpdatedAt time.Time  `json:"updated_at"`
}

type CreateSourceInput struct {
        ProjectID string     `json:"project_id"`
        Name      string     `json:"name"`
        Type      SourceType `json:"type"`
        PositionX float64    `json:"position_x"`
        PositionY float64    `json:"position_y"`
        PositionZ float64    `json:"position_z"`
        AimX      float64    `json:"aim_x"`
        AimY      float64    `json:"aim_y"`
        AimZ      float64    `json:"aim_z"`
        GainDB    float64    `json:"gain_db"`
        GroupID   *string    `json:"group_id"`
}

type UpdateSourceInput struct {
        Name      *string    `json:"name,omitempty"`
        Type      *SourceType `json:"type,omitempty"`
        PositionX *float64   `json:"position_x,omitempty"`
        PositionY *float64   `json:"position_y,omitempty"`
        PositionZ *float64   `json:"position_z,omitempty"`
        AimX      *float64   `json:"aim_x,omitempty"`
        AimY      *float64   `json:"aim_y,omitempty"`
        AimZ      *float64   `json:"aim_z,omitempty"`
        GainDB    *float64   `json:"gain_db,omitempty"`
        GroupID   *string    `json:"group_id,omitempty"`
}

func NewSource(input CreateSourceInput) *Source {
        now := time.Now()
        
        source := &Source{
                ID:        uuid.New().String(),
                ProjectID: input.ProjectID,
                Name:      input.Name,
                Type:      input.Type,
                PositionX: input.PositionX,
                PositionY: input.PositionY,
                PositionZ: input.PositionZ,
                AimX:      input.AimX,
                AimY:      input.AimY,
                AimZ:      input.AimZ,
                GainDB:    input.GainDB,
                GroupID:   input.GroupID,
                CreatedAt: now,
                UpdatedAt: now,
        }
        
        return source
}

func (s *Source) Update(input UpdateSourceInput) {
        if input.Name != nil {
                s.Name = *input.Name
        }
        if input.Type != nil {
                s.Type = *input.Type
        }
        if input.PositionX != nil {
                s.PositionX = *input.PositionX
        }
        if input.PositionY != nil {
                s.PositionY = *input.PositionY
        }
        if input.PositionZ != nil {
                s.PositionZ = *input.PositionZ
        }
        if input.AimX != nil {
                s.AimX = *input.AimX
        }
        if input.AimY != nil {
                s.AimY = *input.AimY
        }
        if input.AimZ != nil {
                s.AimZ = *input.AimZ
        }
        if input.GainDB != nil {
                s.GainDB = *input.GainDB
        }
        if input.GroupID != nil {
                s.GroupID = input.GroupID
        }
        s.UpdatedAt = time.Now()
}

func (st SourceType) String() string {
        return string(st)
}

func (st SourceType) DisplayName() string {
        switch st {
        case SourceTypeSpeaker:
                return "Speaker"
        case SourceTypePointSource:
                return "Point Source"
        case SourceTypeLineArray:
                return "Line Array"
        default:
                return string(st)
        }
}

func ValidSourceTypes() []SourceType {
        return []SourceType{
                SourceTypeSpeaker,
                SourceTypePointSource,
                SourceTypeLineArray,
        }
}
