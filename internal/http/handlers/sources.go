package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/domain"
)

type SourcesHandler struct {
	app *app.App
}

func NewSourcesHandler(app *app.App) *SourcesHandler {
	return &SourcesHandler{app: app}
}

func (h *SourcesHandler) CreateSource(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	input := domain.CreateSourceInput{
		ProjectID: projectID,
		Name:      r.FormValue("name"),
		Type:      domain.SourceType(r.FormValue("type")),
		PositionX: parseFloatOrDefault(r.FormValue("position_x"), 0),
		PositionY: parseFloatOrDefault(r.FormValue("position_y"), 0),
		PositionZ: parseFloatOrDefault(r.FormValue("position_z"), 0),
		GainDB:    parseFloatOrDefault(r.FormValue("gain_db"), 0),
	}

	_, err := h.app.Services.SourceService.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}

func (h *SourcesHandler) DeleteSource(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	sourceID := chi.URLParam(r, "sourceID")
	if projectID == "" || sourceID == "" {
		http.Error(w, "Project ID and Source ID are required", http.StatusBadRequest)
		return
	}

	if err := h.app.Services.SourceService.Delete(r.Context(), sourceID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}
