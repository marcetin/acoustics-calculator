package repo

import (
	"database/sql"
)

type Repositories struct {
	Project *ProjectRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Project: NewProjectRepository(db),
	}
}
