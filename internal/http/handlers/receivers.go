package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/domain"
)

type ReceiversHandler struct {
	app *app.App
}

func NewReceiversHandler(app *app.App) *ReceiversHandler {
	return &ReceiversHandler{app: app}
}

func (h *ReceiversHandler) CreateReceiver(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	input := domain.CreateReceiverInput{
		ProjectID: projectID,
		Name:      r.FormValue("name"),
		Type:      domain.ReceiverType(r.FormValue("type")),
		PositionX: parseFloatOrDefault(r.FormValue("position_x"), 0),
		PositionY: parseFloatOrDefault(r.FormValue("position_y"), 0),
		PositionZ: parseFloatOrDefault(r.FormValue("position_z"), 0),
		Weight:    parseFloatOrDefault(r.FormValue("weight"), 1.0),
	}

	_, err := h.app.Services.ReceiverService.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}

func (h *ReceiversHandler) DeleteReceiver(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	receiverID := chi.URLParam(r, "receiverID")
	if projectID == "" || receiverID == "" {
		http.Error(w, "Project ID and Receiver ID are required", http.StatusBadRequest)
		return
	}

	if err := h.app.Services.ReceiverService.Delete(r.Context(), receiverID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}
