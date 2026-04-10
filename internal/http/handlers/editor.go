package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/domain"
)

type EditorHandler struct {
	app *app.App
}

func NewEditorHandler(app *app.App) *EditorHandler {
	return &EditorHandler{app: app}
}

func (h *EditorHandler) ShowEditor(w http.ResponseWriter, r *http.Request) {
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

	data := map[string]interface{}{
		"Project":   project,
		"Geometry":  geometry,
		"Sources":   sources,
		"Receivers": receivers,
	}

	if err := h.app.Templates.ExecuteTemplate(w, "room_editor.html", data); err != nil {
		http.Error(w, "Failed to render editor", http.StatusInternalServerError)
	}
}

func (h *EditorHandler) MoveSource(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	sourceID := chi.URLParam(r, "sourceID")
	if projectID == "" || sourceID == "" {
		http.Error(w, "Project ID and Source ID are required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	x, _ := strconv.ParseFloat(r.FormValue("x"), 64)
	y, _ := strconv.ParseFloat(r.FormValue("y"), 64)

	input := domain.UpdateSourceInput{
		PositionX: &x,
		PositionY: &y,
	}

	if _, err := h.app.Services.SourceService.Update(r.Context(), sourceID, input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EditorHandler) MoveReceiver(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	receiverID := chi.URLParam(r, "receiverID")
	if projectID == "" || receiverID == "" {
		http.Error(w, "Project ID and Receiver ID are required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	x, _ := strconv.ParseFloat(r.FormValue("x"), 64)
	y, _ := strconv.ParseFloat(r.FormValue("y"), 64)

	input := domain.UpdateReceiverInput{
		PositionX: &x,
		PositionY: &y,
	}

	if _, err := h.app.Services.ReceiverService.Update(r.Context(), receiverID, input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
