package domain

import "time"

type ProjectReport struct {
	Project          *ProjectReportProject     `json:"project"`
	Geometry         *ProjectReportGeometry    `json:"geometry,omitempty"`
	Surfaces         []*ProjectReportSurface   `json:"surfaces,omitempty"`
	Materials        []*ProjectReportMaterial  `json:"materials,omitempty"`
	Sources          []*ProjectReportSource    `json:"sources,omitempty"`
	Receivers        []*ProjectReportReceiver  `json:"receivers,omitempty"`
	Constraints      *ProjectReportConstraints `json:"constraints,omitempty"`
	LatestAnalysis   *ProjectReportAnalysis    `json:"latest_analysis,omitempty"`
	LatestPlacements *ProjectReportPlacements  `json:"latest_placements,omitempty"`
	Metadata         *ProjectReportMetadata    `json:"metadata"`
}

type ProjectReportProject struct {
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

type ProjectReportGeometry struct {
	GeometryType GeometryType `json:"geometry_type"`
	WidthM       float64      `json:"width_m"`
	LengthM      float64      `json:"length_m"`
	HeightM      float64      `json:"height_m"`
	VolumeM3     float64      `json:"volume_m3"`
}

type ProjectReportSurface struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Kind         SurfaceKind `json:"kind"`
	Orientation  string      `json:"orientation"`
	AreaM2       float64     `json:"area_m2"`
	Mountable    bool        `json:"mountable"`
	MaterialID   *string     `json:"material_id,omitempty"`
	MaterialName *string     `json:"material_name,omitempty"`
}

type ProjectReportMaterial struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	IsPreset        bool            `json:"is_preset"`
	AbsorptionCoeff map[int]float64 `json:"absorption_coeff,omitempty"`
	ScatteringCoeff map[int]float64 `json:"scattering_coeff,omitempty"`
}

type ProjectReportSource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Z         float64   `json:"z"`
	AimX      float64   `json:"aim_x"`
	AimY      float64   `json:"aim_y"`
	AimZ      float64   `json:"aim_z"`
	GainDB    float64   `json:"gain_db"`
	GroupID   *string   `json:"group_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectReportReceiver struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Z         float64   `json:"z"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectReportConstraints struct {
	MaxPanelDepthM         float64       `json:"max_panel_depth_m"`
	SymmetryRequired       bool          `json:"symmetry_required"`
	BudgetClass            string        `json:"budget_class"`
	PreferredTreatmentMode string        `json:"preferred_treatment_mode"`
	AllowedSurfaceKinds    []SurfaceKind `json:"allowed_surface_kinds,omitempty"`
	ForbiddenSurfaces      []string      `json:"forbidden_surfaces,omitempty"`
}

type ProjectReportAnalysis struct {
	ID           string           `json:"id"`
	ProjectID    string           `json:"project_id"`
	Status       string           `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	InputsHash   string           `json:"inputs_hash,omitempty"`
	SpeedOfSound float64          `json:"speed_of_sound"`
	Summary      *AnalysisSummary `json:"summary,omitempty"`
	Warnings     []string         `json:"warnings,omitempty"`
	Modes        []Mode           `json:"modes,omitempty"`
	RT           RTResults        `json:"rt,omitempty"`
	Reflections  []ReflectionPath `json:"reflections,omitempty"`
	IsStale      bool             `json:"is_stale"`
}

type ProjectReportPlacements struct {
	AnalysisRunID *string               `json:"analysis_run_id,omitempty"`
	Summary       *PlacementSummary     `json:"summary,omitempty"`
	Warnings      []string              `json:"warnings,omitempty"`
	Candidates    []*PlacementCandidate `json:"candidates,omitempty"`
	IsStale       bool                  `json:"is_stale"`
}

type ProjectReportMetadata struct {
	GeneratedAt       time.Time `json:"generated_at"`
	ReportVersion     string    `json:"report_version"`
	SupportedScope    string    `json:"supported_scope"`
	Assumptions       []string  `json:"assumptions"`
	IsPublicProfile   bool      `json:"is_public_profile"`
	PublicProvisional bool      `json:"public_provisional"`
}

type PlacementCSVRow struct {
	ProjectID       string
	AnalysisRunID   string
	PlacementID     string
	SurfaceID       string
	SurfaceName     string
	Decision        string
	Score           float64
	Confidence      float64
	DiffuserType    string
	Orientation     string
	CenterX         float64
	CenterY         float64
	CenterZ         float64
	PanelCount      int
	CoverageAreaM2  float64
	TargetBands     string
	ReasonsSummary  string
	WarningsSummary string
	CreatedAt       time.Time
}
