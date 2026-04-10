package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"acoustics-calculator/internal/app"
)

type ReportsHandler struct {
	app *app.App
}

func NewReportsHandler(app *app.App) *ReportsHandler {
	return &ReportsHandler{app: app}
}

func (h *ReportsHandler) ExportJSON(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	report, err := h.app.Services.ReportService.BuildProjectReport(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-report.json", report.Project.Name))

	if err := json.NewEncoder(w).Encode(report); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ReportsHandler) ExportPlacementsCSV(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	rows, err := h.app.Services.ReportService.BuildPlacementsCSV(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	project, err := h.app.Services.ProjectService.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-placements.csv", project.Name))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{
		"ProjectID",
		"AnalysisRunID",
		"PlacementID",
		"SurfaceID",
		"SurfaceName",
		"Decision",
		"Score",
		"Confidence",
		"DiffuserType",
		"Orientation",
		"CenterX",
		"CenterY",
		"CenterZ",
		"PanelCount",
		"CoverageAreaM2",
		"TargetBands",
		"ReasonsSummary",
		"WarningsSummary",
		"CreatedAt",
	}
	if err := writer.Write(header); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, row := range rows {
		record := []string{
			row.ProjectID,
			row.AnalysisRunID,
			row.PlacementID,
			row.SurfaceID,
			row.SurfaceName,
			row.Decision,
			fmt.Sprintf("%.4f", row.Score),
			fmt.Sprintf("%.4f", row.Confidence),
			row.DiffuserType,
			row.Orientation,
			fmt.Sprintf("%.4f", row.CenterX),
			fmt.Sprintf("%.4f", row.CenterY),
			fmt.Sprintf("%.4f", row.CenterZ),
			strconv.Itoa(row.PanelCount),
			fmt.Sprintf("%.4f", row.CoverageAreaM2),
			row.TargetBands,
			row.ReasonsSummary,
			row.WarningsSummary,
			row.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *ReportsHandler) ShowReport(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	report, err := h.app.Services.ReportService.GetPrintableReportData(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Report": report,
	}

	if err := h.app.Templates.ExecuteTemplate(w, "project_report.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
