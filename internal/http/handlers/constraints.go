package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
	"acoustics-calculator/internal/domain"
)

type ConstraintsHandler struct {
	app *app.App
}

func NewConstraintsHandler(app *app.App) *ConstraintsHandler {
	return &ConstraintsHandler{app: app}
}

func (h *ConstraintsHandler) SaveConstraints(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	maxPanelDepth := parseFloatOrDefault(r.FormValue("max_panel_depth_m"), 0.3)
	symmetryRequired := r.FormValue("symmetry_required") == "on"
	treatmentMode := domain.TreatmentMode(r.FormValue("preferred_treatment_mode"))
	budgetClass := domain.BudgetClass(r.FormValue("budget_class"))

	input := domain.UpdateConstraintInput{
		MaxPanelDepthM:         &maxPanelDepth,
		SymmetryRequired:       &symmetryRequired,
		PreferredTreatmentMode: &treatmentMode,
		BudgetClass:            &budgetClass,
	}

	_, err := h.app.Services.ConstraintService.Save(r.Context(), projectID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/projects/"+projectID, http.StatusSeeOther)
}
