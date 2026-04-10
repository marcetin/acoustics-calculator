package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
)

type JobsHandler struct {
	app *app.App
}

func NewJobsHandler(app *app.App) *JobsHandler {
	return &JobsHandler{app: app}
}

func (h *JobsHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	jobID := chi.URLParam(r, "jobID")

	if projectID == "" || jobID == "" {
		http.Error(w, "Project ID and Job ID are required", http.StatusBadRequest)
		return
	}

	job, err := h.app.Services.JobService.GetJob(r.Context(), jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status":"` + string(job.Status) + `","phase":"` + job.Phase + `","progress":` + string(rune(job.ProgressPercent)) + `,"message":"` + job.Message + `"}`))
}
