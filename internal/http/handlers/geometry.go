package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/domain"
)

type GeometryHandler struct {
	app *app.App
}

func NewGeometryHandler(app *app.App) *GeometryHandler {
	return &GeometryHandler{app: app}
}

func (h *GeometryHandler) ShowGeometry(w http.ResponseWriter, r *http.Request) {
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
		"Title":    project.Name + " - Geometry",
		"Project":  project,
		"Geometry": geometry,
		"Tab":      "geometry",
	}

	if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *GeometryHandler) SaveGeometry(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	input := domain.CreateGeometryInput{
		ProjectID:    projectID,
		GeometryType: domain.GeometryType(r.FormValue("geometry_type")),
		Width:        parseFloatOrDefault(r.FormValue("width"), 0),
		Length:       parseFloatOrDefault(r.FormValue("length"), 0),
		Height:       parseFloatOrDefault(r.FormValue("height"), 0),
	}

	_, err := h.app.Services.GeometryService.UpsertShoeboxGeometry(r.Context(), input)
	if err != nil {
		data := map[string]interface{}{
			"Title":     "Project - Geometry",
			"ProjectID": projectID,
			"Form":      input,
			"Error":     err.Error(),
			"Geometry":  nil,
		}

		w.WriteHeader(http.StatusBadRequest)
		if err := h.app.Templates.ExecuteTemplate(w, "base.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}
