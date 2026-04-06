package app

import (
	"database/sql"
	"html/template"

	"acoustics-calculator/internal/repo"
	"acoustics-calculator/internal/service"
)

type App struct {
	DB        *sql.DB
	Repos     *repo.Repositories
	Services  *Services
	Templates *template.Template
	Config    *Config
}

type Services struct {
	ProjectService *service.ProjectService
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
		ProjectService: service.NewProjectService(a.Repos.Project),
	}

	if err := a.loadTemplates(); err != nil {
		return err
	}

	return nil
}

func (a *App) loadTemplates() error {
	templates := template.Must(template.ParseGlob("templates/layouts/*.html"))
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
