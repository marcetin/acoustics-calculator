package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SurfaceKind string

const (
	SurfaceKindWall      SurfaceKind = "WALL"
	SurfaceKindCeiling   SurfaceKind = "CEILING"
	SurfaceKindFloor     SurfaceKind = "FLOOR"
	SurfaceKindObjectFace SurfaceKind = "OBJECT_FACE"
)

type Surface struct {
	ID             string       `json:"id"`
	ProjectID      string       `json:"project_id"`
	GeometryID     string       `json:"geometry_id"`
	Name           string       `json:"name"`
	Kind           SurfaceKind  `json:"kind"`
	Polygon3DJSON  string       `json:"polygon_3d_json"`
	NormalX        float64      `json:"normal_x"`
	NormalY        float64      `json:"normal_y"`
	NormalZ        float64      `json:"normal_z"`
	AreaM2         float64      `json:"area_m2"`
	MaterialID     *string      `json:"material_id"`
	IsMountable    bool         `json:"is_mountable"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type CreateSurfaceInput struct {
	ProjectID     string       `json:"project_id"`
	GeometryID    string       `json:"geometry_id"`
	Name          string       `json:"name"`
	Kind          SurfaceKind  `json:"kind"`
	Polygon3D     []float64    `json:"polygon_3d"`
	NormalX       float64      `json:"normal_x"`
	NormalY       float64      `json:"normal_y"`
	NormalZ       float64      `json:"normal_z"`
	AreaM2        float64      `json:"area_m2"`
	MaterialID    *string      `json:"material_id"`
	IsMountable   bool         `json:"is_mountable"`
}

type UpdateSurfaceInput struct {
	Name        *string      `json:"name,omitempty"`
	MaterialID  *string      `json:"material_id,omitempty"`
	IsMountable *bool        `json:"is_mountable,omitempty"`
}

func NewSurface(input CreateSurfaceInput) (*Surface, error) {
	now := time.Now()
	
	polygonJSON, err := json.Marshal(input.Polygon3D)
	if err != nil {
		return nil, err
	}
	
	surface := &Surface{
		ID:            uuid.New().String(),
		ProjectID:     input.ProjectID,
		GeometryID:    input.GeometryID,
		Name:          input.Name,
		Kind:          input.Kind,
		Polygon3DJSON: string(polygonJSON),
		NormalX:       input.NormalX,
		NormalY:       input.NormalY,
		NormalZ:       input.NormalZ,
		AreaM2:        input.AreaM2,
		MaterialID:    input.MaterialID,
		IsMountable:   input.IsMountable,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	return surface, nil
}

func (s *Surface) Update(input UpdateSurfaceInput) {
	if input.Name != nil {
		s.Name = *input.Name
	}
	if input.MaterialID != nil {
		s.MaterialID = input.MaterialID
	}
	if input.IsMountable != nil {
		s.IsMountable = *input.IsMountable
	}
	s.UpdatedAt = time.Now()
}

func (s *Surface) GetPolygon3D() ([]float64, error) {
	var polygon []float64
	if s.Polygon3DJSON == "" {
		return polygon, nil
	}
	err := json.Unmarshal([]byte(s.Polygon3DJSON), &polygon)
	return polygon, err
}

func (sk SurfaceKind) String() string {
	return string(sk)
}

func (sk SurfaceKind) DisplayName() string {
	switch sk {
	case SurfaceKindWall:
		return "Wall"
	case SurfaceKindCeiling:
		return "Ceiling"
	case SurfaceKindFloor:
		return "Floor"
	case SurfaceKindObjectFace:
		return "Object Face"
	default:
		return string(sk)
	}
}

func ValidSurfaceKinds() []SurfaceKind {
	return []SurfaceKind{
		SurfaceKindWall,
		SurfaceKindCeiling,
		SurfaceKindFloor,
		SurfaceKindObjectFace,
	}
}
