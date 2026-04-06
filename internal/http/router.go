package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/http/handlers"
)

func New(application *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Initialize handlers
	projectHandler := handlers.NewProjectHandler(application)

	// Routes
	r.Get("/", projectHandler.ListProjects)
	r.Get("/dashboard", projectHandler.ListProjects)

	// Project routes
	r.Get("/projects/new", projectHandler.ShowNewProjectForm)
	r.Post("/projects", projectHandler.CreateProject)
	r.Get("/projects/{id}", projectHandler.ShowProject)

	// HTMX tab routes
	r.Get("/hx/projects/{id}/tabs/{tab}", projectHandler.ShowTab)

	return r
}
