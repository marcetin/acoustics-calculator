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
	geometryHandler := handlers.NewGeometryHandler(application)
	surfacesHandler := handlers.NewSurfacesHandler(application)
	sourcesHandler := handlers.NewSourcesHandler(application)
	receiversHandler := handlers.NewReceiversHandler(application)
	constraintsHandler := handlers.NewConstraintsHandler(application)
	editorHandler := handlers.NewEditorHandler(application)
	analysisHandler := handlers.NewAnalysisHandler(application)
	placementsHandler := handlers.NewPlacementsHandler(application)

	// Routes
	r.Get("/", projectHandler.ListProjects)
	r.Get("/dashboard", projectHandler.ListProjects)

	// Project routes
	r.Get("/projects/new", projectHandler.ShowNewProjectForm)
	r.Post("/projects", projectHandler.CreateProject)
	r.Get("/projects/{id}", projectHandler.ShowProject)

	// Geometry routes
	r.Post("/projects/{id}/geometry", geometryHandler.SaveGeometry)

	// Surfaces routes
	r.Post("/projects/{id}/surfaces/{surfaceID}", surfacesHandler.UpdateSurface)

	// Sources routes
	r.Post("/projects/{id}/sources", sourcesHandler.CreateSource)
	r.Post("/projects/{id}/sources/{sourceID}/delete", sourcesHandler.DeleteSource)

	// Receivers routes
	r.Post("/projects/{id}/receivers", receiversHandler.CreateReceiver)
	r.Post("/projects/{id}/receivers/{receiverID}/delete", receiversHandler.DeleteReceiver)

	// Constraints routes
	r.Post("/projects/{id}/constraints", constraintsHandler.SaveConstraints)

	// Editor routes
	r.Get("/hx/projects/{id}/editor", editorHandler.ShowEditor)
	r.Post("/hx/projects/{id}/editor/sources/{sourceID}/move", editorHandler.MoveSource)
	r.Post("/hx/projects/{id}/editor/receivers/{receiverID}/move", editorHandler.MoveReceiver)

	// Analysis routes
	r.Post("/projects/{id}/analysis/run", analysisHandler.RunAnalysis)

	// Placements routes
	r.Post("/projects/{id}/placements/generate", placementsHandler.GeneratePlacements)

	// HTMX tab routes
	r.Get("/hx/projects/{id}/tabs/{tab}", projectHandler.ShowTab)

	return r
}
