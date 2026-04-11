package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"acoustics-calculator/internal/scene"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	
	HandleHealth(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	if !response.OK {
		t.Error("Expected OK to be true")
	}
}

func TestHandleGetExample(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/scene/example", nil)
	w := httptest.NewRecorder()
	
	HandleGetExample(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response SuccessResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	if !response.OK {
		t.Error("Expected OK to be true")
	}
	
	// Verify the data is a valid SceneConfig
	var config scene.SceneConfig
	dataBytes, _ := json.Marshal(response.Data)
	if err := json.Unmarshal(dataBytes, &config); err != nil {
		t.Fatal(err)
	}
	
	if config.Room.WidthM <= 0 {
		t.Error("Expected positive room width")
	}
}

func TestHandlePreview_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/scene/preview", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()
	
	HandlePreview(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandlePreview_InvalidScene(t *testing.T) {
	invalidConfig := scene.SceneConfig{
		Room: scene.RoomConfig{
			WidthM: -1, // Invalid
		},
	}
	body, _ := json.Marshal(invalidConfig)
	
	req := httptest.NewRequest("POST", "/api/scene/preview", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	HandlePreview(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	
	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	if response.OK {
		t.Error("Expected OK to be false for invalid scene")
	}
	
	if response.Error.Code != "validation_error" {
		t.Errorf("Expected validation_error code, got %s", response.Error.Code)
	}
}

func TestHandlePreview_ValidScene(t *testing.T) {
	validConfig := scene.GetExample()
	body, _ := json.Marshal(validConfig)
	
	req := httptest.NewRequest("POST", "/api/scene/preview", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	
	HandlePreview(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response SuccessResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	
	if !response.OK {
		t.Error("Expected OK to be true")
	}
}
