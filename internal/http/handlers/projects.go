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
		"Title":     "New Project",
		"SpaceTypes": domain.ValidSpaceTypes(),
		"GoalTypes":  domain.ValidGoalTypes(),
		"Units":     domain.ValidUnits(),
	}

	if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
			"Title":       "New Project",
			"SpaceTypes":  domain.ValidSpaceTypes(),
			"GoalTypes":   domain.ValidGoalTypes(),
			"Units":       domain.ValidUnits(),
			"Form":        input,
			"Error":       err.Error(),
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

	data := map[string]interface{}{
		"Title":   project.Name,
		"Project": project,
		"Tab":     "summary",
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

	data := map[string]interface{}{
		"Project": project,
		"Tab":     tabName,
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
