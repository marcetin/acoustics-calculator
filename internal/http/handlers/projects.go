package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/domain"
)

type ProjectHandler struct {
	app *app.App
}

func NewProjectHandler(app *app.App) *ProjectHandler {
	return &ProjectHandler{app: app}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.app.Services.ProjectService.ListProjects(r.Context())
	if err != nil {
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":    "Dashboard",
		"Projects": projects,
	}

	if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProjectHandler) ShowNewProjectForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":      "New Project",
		"SpaceTypes": domain.ValidSpaceTypes(),
		"GoalTypes":  domain.ValidGoalTypes(),
		"Units":      domain.ValidUnits(),
	}

	if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProjectHandler) CreateDemoProject(w http.ResponseWriter, r *http.Request) {
	project, err := h.app.Services.ProjectService.CreateDemoProject(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects/"+project.ID, http.StatusSeeOther)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	input := domain.CreateProjectInput{
		Name:            r.FormValue("name"),
		SpaceType:       domain.SpaceType(r.FormValue("space_type")),
		GoalType:        domain.GoalType(r.FormValue("goal_type")),
		Units:           domain.Units(r.FormValue("units")),
		TemperatureC:    parseFloatOrDefault(r.FormValue("temperature_c"), 20.0),
		HumidityPercent: parseFloatOrDefault(r.FormValue("humidity_percent"), 50.0),
		SpeedOfSound:    parseFloatOrDefault(r.FormValue("speed_of_sound"), 343.0),
	}

	project, err := h.app.Services.ProjectService.CreateProject(r.Context(), input)
	if err != nil {
		data := map[string]interface{}{
			"Title":      "New Project",
			"SpaceTypes": domain.ValidSpaceTypes(),
			"GoalTypes":  domain.ValidGoalTypes(),
			"Units":      domain.ValidUnits(),
			"Form":       input,
			"Error":      err.Error(),
		}

		w.WriteHeader(http.StatusBadRequest)
		if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/projects/"+project.ID, http.StatusSeeOther)
}

func (h *ProjectHandler) ShowProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	project, err := h.app.Services.ProjectService.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	geometry, err := h.app.Services.GeometryService.GetByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to load geometry", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":    project.Name,
		"Project":  project,
		"Geometry": geometry,
		"Tab":      "summary",
	}

	if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ProjectHandler) ShowTab(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tabName := chi.URLParam(r, "tab")
	if projectID == "" || tabName == "" {
		http.Error(w, "Project ID and tab name are required", http.StatusBadRequest)
		return
	}

	project, err := h.app.Services.ProjectService.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	geometry, err := h.app.Services.GeometryService.GetByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to load geometry", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Project":  project,
		"Geometry": geometry,
		"Tab":      tabName,
	}

	if tabName == "surfaces" {
		surfaces, err := h.app.Services.MaterialService.ListProjectSurfaces(r.Context(), projectID)
		if err != nil {
			http.Error(w, "Failed to load surfaces", http.StatusInternalServerError)
			return
		}
		materials, err := h.app.Services.MaterialService.ListPresets(r.Context())
		if err != nil {
			http.Error(w, "Failed to load materials", http.StatusInternalServerError)
			return
		}
		data["Surfaces"] = surfaces
		data["Materials"] = materials
	}

	if tabName == "sources" {
		sources, err := h.app.Services.SourceService.ListByProject(r.Context(), projectID)
		if err != nil {
			http.Error(w, "Failed to load sources", http.StatusInternalServerError)
			return
		}
		data["Sources"] = sources
	}

	if tabName == "receivers" {
		receivers, err := h.app.Services.ReceiverService.ListByProject(r.Context(), projectID)
		if err != nil {
			http.Error(w, "Failed to load receivers", http.StatusInternalServerError)
			return
		}
		data["Receivers"] = receivers
	}

	if tabName == "summary" {
		surfaces, err := h.app.Services.MaterialService.ListProjectSurfaces(r.Context(), projectID)
		if err != nil {
			http.Error(w, "Failed to load surfaces", http.StatusInternalServerError)
			return
		}

		sources, err := h.app.Services.SourceService.ListByProject(r.Context(), projectID)
		if err != nil {
			http.Error(w, "Failed to load sources", http.StatusInternalServerError)
			return
		}

		receivers, err := h.app.Services.ReceiverService.ListByProject(r.Context(), projectID)
		if err != nil {
			http.Error(w, "Failed to load receivers", http.StatusInternalServerError)
			return
		}

		constraints, err := h.app.Services.ConstraintService.GetOrCreateDefault(r.Context(), projectID)
		if err == nil {
			data["Constraints"] = constraints
		}

		materialsAssigned := 0
		for _, surface := range surfaces {
			if surface.MaterialID != nil {
				materialsAssigned++
			}
		}

		readyForAnalysis := geometry != nil && len(sources) > 0 && len(receivers) > 0

		data["Surfaces"] = surfaces
		data["Sources"] = sources
		data["Receivers"] = receivers
		data["MaterialsAssigned"] = materialsAssigned
		data["ReadyForAnalysis"] = readyForAnalysis
	}

	if tabName == "analysis" {
		latestRun, err := h.app.Services.AnalysisService.GetLatestByProject(r.Context(), projectID)
		if err == nil && latestRun != nil {
			metrics, err := latestRun.GetMetrics()
			if err == nil {
				data["AnalysisRun"] = latestRun
				data["AnalysisMetrics"] = metrics

				isStale, err := h.app.Services.AnalysisService.IsAnalysisStale(r.Context(), projectID, latestRun)
				if err == nil {
					data["AnalysisStale"] = isStale
				}
			}
		}

		data["Geometry"] = geometry
		data["ReadyForAnalysis"] = geometry != nil

		readiness, err := h.app.Services.AnalysisService.CanRunAnalysis(r.Context(), projectID)
		if err == nil {
			data["AnalysisReadiness"] = readiness
		}

		activeJob, err := h.app.Services.JobService.GetLatestActiveByProject(r.Context(), projectID)
		if err == nil && activeJob != nil && activeJob.Kind == domain.JobKindAnalysis {
			data["ActiveAnalysisJob"] = activeJob
		}
	}

	if tabName == "placements" {
		readiness, err := h.app.Services.PlacementService.CanGeneratePlacements(r.Context(), projectID)
		if err == nil {
			data["PlacementReadiness"] = readiness
		}

		placements, err := h.app.Services.PlacementService.GetLatestPlacements(r.Context(), projectID)
		if err == nil {
			data["Placements"] = placements

			isStale, err := h.app.Services.PlacementService.ArePlacementsStale(r.Context(), projectID, placements)
			if err == nil {
				data["PlacementsStale"] = isStale
			}
		}

		summary, err := h.app.Services.PlacementService.SummarizePlacements(r.Context(), projectID)
		if err == nil {
			data["PlacementSummary"] = summary
		}

		latestRun, err := h.app.Services.AnalysisService.GetLatestByProject(r.Context(), projectID)
		if err == nil && latestRun != nil {
			data["LatestAnalysisRun"] = latestRun
		}

		activeJob, err := h.app.Services.JobService.GetLatestActiveByProject(r.Context(), projectID)
		if err == nil && activeJob != nil && activeJob.Kind == domain.JobKindPlacement {
			data["ActivePlacementJob"] = activeJob
		}
	}

	templateName := "tab_" + tabName + ".html"
	if err := h.app.Templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "Tab not found", http.StatusNotFound)
	}
}

func parseFloatOrDefault(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	if value, err := strconv.ParseFloat(s, 64); err == nil {
		return value
	}
	return defaultValue
}
