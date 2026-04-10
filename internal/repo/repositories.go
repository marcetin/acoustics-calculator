package repo

import (
	"database/sql"
)

type Repositories struct {
	Project    *ProjectRepository
	Geometry   *GeometryRepository
	Surface    *SurfaceRepository
	Material   *MaterialRepository
	Source     *SourceRepository
	Receiver   *ReceiverRepository
	Constraint *ConstraintRepository
	Analysis   *AnalysisRepository
	Diffuser   *DiffuserRepository
	Placement  *PlacementRepository
	Job        *JobRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Project:    NewProjectRepository(db),
		Geometry:   NewGeometryRepository(db),
		Surface:    NewSurfaceRepository(db),
		Material:   NewMaterialRepository(db),
		Source:     NewSourceRepository(db),
		Receiver:   NewReceiverRepository(db),
		Constraint: NewConstraintRepository(db),
		Analysis:   NewAnalysisRepository(db),
		Diffuser:   NewDiffuserRepository(db),
		Placement:  NewPlacementRepository(db),
		Job:        NewJobRepository(db),
	}
}
