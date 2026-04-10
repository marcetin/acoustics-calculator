package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReceiverType string

const (
	ReceiverTypeListener    ReceiverType = "LISTENER"
	ReceiverTypeMic         ReceiverType = "MIC"
	ReceiverTypeAudienceZone ReceiverType = "AUDIENCE_ZONE"
)

type Receiver struct {
	ID        string       `json:"id"`
	ProjectID string       `json:"project_id"`
	Name      string       `json:"name"`
	Type      ReceiverType `json:"type"`
	PositionX float64      `json:"position_x"`
	PositionY float64      `json:"position_y"`
	PositionZ float64      `json:"position_z"`
	Weight    float64      `json:"weight"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateReceiverInput struct {
	ProjectID string       `json:"project_id"`
	Name      string       `json:"name"`
	Type      ReceiverType `json:"type"`
	PositionX float64      `json:"position_x"`
	PositionY float64      `json:"position_y"`
	PositionZ float64      `json:"position_z"`
	Weight    float64      `json:"weight"`
}

type UpdateReceiverInput struct {
	Name      *string       `json:"name,omitempty"`
	Type      *ReceiverType `json:"type,omitempty"`
	PositionX *float64      `json:"position_x,omitempty"`
	PositionY *float64      `json:"position_y,omitempty"`
	PositionZ *float64      `json:"position_z,omitempty"`
	Weight    *float64      `json:"weight,omitempty"`
}

func NewReceiver(input CreateReceiverInput) *Receiver {
	now := time.Now()
	
	receiver := &Receiver{
		ID:        uuid.New().String(),
		ProjectID: input.ProjectID,
		Name:      input.Name,
		Type:      input.Type,
		PositionX: input.PositionX,
		PositionY: input.PositionY,
		PositionZ: input.PositionZ,
		Weight:    input.Weight,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	return receiver
}

func (r *Receiver) Update(input UpdateReceiverInput) {
	if input.Name != nil {
		r.Name = *input.Name
	}
	if input.Type != nil {
		r.Type = *input.Type
	}
	if input.PositionX != nil {
		r.PositionX = *input.PositionX
	}
	if input.PositionY != nil {
		r.PositionY = *input.PositionY
	}
	if input.PositionZ != nil {
		r.PositionZ = *input.PositionZ
	}
	if input.Weight != nil {
		r.Weight = *input.Weight
	}
	r.UpdatedAt = time.Now()
}

func (rt ReceiverType) String() string {
	return string(rt)
}

func (rt ReceiverType) DisplayName() string {
	switch rt {
	case ReceiverTypeListener:
		return "Listener"
	case ReceiverTypeMic:
		return "Microphone"
	case ReceiverTypeAudienceZone:
		return "Audience Zone"
	default:
		return string(rt)
	}
}

func ValidReceiverTypes() []ReceiverType {
	return []ReceiverType{
		ReceiverTypeListener,
		ReceiverTypeMic,
		ReceiverTypeAudienceZone,
	}
}
