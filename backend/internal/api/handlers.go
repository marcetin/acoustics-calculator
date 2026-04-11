package api

import (
	"encoding/json"
	"net/http"

	"acoustics-calculator/internal/materials"
	"acoustics-calculator/internal/scene"
	"acoustics-calculator/internal/validation"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{OK: true})
}

func HandleGetExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	example := scene.GetExample()
	json.NewEncoder(w).Encode(SuccessResponse{
		OK:       true,
		Data:     example,
		Warnings: []string{},
	})
}

func HandlePreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var config scene.SceneConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			OK: false,
			Error: ErrorDetail{
				Code:    "invalid_json",
				Message: "Invalid JSON in request body",
			},
		})
		return
	}

	// Validate
	materialCatalog := materials.GetCatalog()
	errors := validation.ValidateScene(config, materialCatalog)
	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			OK: false,
			Error: ErrorDetail{
				Code:    "validation_error",
				Message: "Invalid scene configuration",
				Fields:  errors,
			},
		})
		return
	}

	// Derive scene
	derived, warnings := scene.DeriveScene(config, materialCatalog)

	response := map[string]interface{}{
		"config":   config,
		"derived":  derived,
		"warnings": warnings,
	}

	json.NewEncoder(w).Encode(SuccessResponse{
		OK:       true,
		Data:     response,
		Warnings: warnings,
	})
}
