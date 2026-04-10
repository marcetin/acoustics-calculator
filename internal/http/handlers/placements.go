package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
)

type PlacementsHandler struct {
	app *app.App
}

func NewPlacementsHandler(app *app.App) *PlacementsHandler {
	return &PlacementsHandler{app: app}
}

func (h *PlacementsHandler) GeneratePlacements(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	_, err := h.app.Services.PlacementService.GeneratePlacements(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}
