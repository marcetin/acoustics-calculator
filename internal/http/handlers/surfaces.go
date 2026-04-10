package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
)

type SurfacesHandler struct {
	app *app.App
}

func NewSurfacesHandler(app *app.App) *SurfacesHandler {
	return &SurfacesHandler{app: app}
}

func (h *SurfacesHandler) UpdateSurface(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	surfaceID := chi.URLParam(r, "surfaceID")
	if projectID == "" || surfaceID == "" {
		http.Error(w, "Project ID and Surface ID are required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	materialID := r.FormValue("material_id")
	var materialIDPtr *string
	if materialID != "" {
		materialIDPtr = &materialID
	}

	if err := h.app.Services.MaterialService.AssignMaterialToSurface(r.Context(), surfaceID, materialIDPtr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
