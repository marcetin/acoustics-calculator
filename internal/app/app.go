package app

import (
	"database/sql"
	"html/template"
	"reflect"

	"acoustics-calculator/internal/repo"
	"acoustics-calculator/internal/service"
)

func firstN(n int, slice interface{}) interface{} {
	if slice == nil {
		return nil
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice
	}

	if n >= v.Len() {
		return slice
	}

	return v.Slice(0, n).Interface()
}

type App struct {
	DB        *sql.DB
	Repos     *repo.Repositories
	Services  *Services
	Templates *template.Template
	Config    *Config
}

type Services struct {
	ProjectService    *service.ProjectService
	GeometryService   *service.GeometryService
	MaterialService   *service.MaterialService
	SourceService     *service.SourceService
	ReceiverService   *service.ReceiverService
	ConstraintService *service.ConstraintService
	AnalysisService   *service.AnalysisService
	PlacementService  *service.PlacementService
}

func New() *App {
	return &App{
		Config: NewConfig(),
	}
}

func (a *App) Bootstrap() error {
	db, err := a.Config.Database.Open()
	if err != nil {
		return err
	}

	a.DB = db
	a.Repos = repo.NewRepositories(db)

	if err := a.runMigrations(); err != nil {
		return err
	}

	a.Services = &Services{
		ProjectService:    service.NewProjectService(a.Repos.Project),
		GeometryService:   service.NewGeometryService(a.Repos.Geometry, a.Repos.Surface),
		MaterialService:   service.NewMaterialService(a.Repos.Material, a.Repos.Surface),
		SourceService:     service.NewSourceService(a.Repos.Source, a.Repos.Geometry),
		ReceiverService:   service.NewReceiverService(a.Repos.Receiver, a.Repos.Geometry),
		ConstraintService: service.NewConstraintService(a.Repos.Constraint),
		AnalysisService:   service.NewAnalysisService(a.Repos.Analysis, a.Repos.Project, a.Repos.Geometry, a.Repos.Surface, a.Repos.Source, a.Repos.Receiver, a.Repos.Material),
		PlacementService:  service.NewPlacementService(a.Repos.Placement, a.Repos.Project, a.Repos.Geometry, a.Repos.Surface, a.Repos.Analysis, a.Repos.Diffuser, a.Repos.Constraint),
	}

	if err := a.loadTemplates(); err != nil {
		return err
	}

	return nil
}

func (a *App) loadTemplates() error {
	funcMap := template.FuncMap{
		"firstN": firstN,
	}
	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/layouts/*.html"))
	templates = template.Must(templates.ParseGlob("templates/pages/*.html"))
	templates = template.Must(templates.ParseGlob("templates/partials/*.html"))

	a.Templates = templates
	return nil
}

func (a *App) runMigrations() error {
	return repo.RunMigrations(a.DB)
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
