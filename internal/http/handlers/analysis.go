package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
)

type AnalysisHandler struct {
	app *app.App
}

func NewAnalysisHandler(app *app.App) *AnalysisHandler {
	return &AnalysisHandler{app: app}
}

func (h *AnalysisHandler) RunAnalysis(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	_, err := h.app.Services.AnalysisService.RunAnalysis(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}
