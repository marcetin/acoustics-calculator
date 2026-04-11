package validation

import (
	"testing"

	"acoustics-calculator/internal/materials"
	"acoustics-calculator/internal/scene"
)

func TestValidateScene_ValidExample(t *testing.T) {
	config := scene.GetExample()
	catalog := materials.GetCatalog()
	
	errors := ValidateScene(config, catalog)
	
	if len(errors) > 0 {
		t.Errorf("Expected no errors for example scene, got %d", len(errors))
	}
}

func TestValidateScene_InvalidRoomDimensions(t *testing.T) {
	config := scene.SceneConfig{
		Room: scene.RoomConfig{
			WidthM:  -1,
			LengthM: 5,
			HeightM: 3,
		},
		Surfaces: scene.SurfaceConfig{
			FrontWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RearWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			LeftWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RightWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			Ceiling:    scene.SurfaceMaterial{MaterialID: "painted_concrete"},
			Floor:      scene.SurfaceMaterial{MaterialID: "carpet"},
		},
		Sources: scene.SourcesConfig{
			Left:  scene.Position3D{X: 1, Y: 1, Z: 1},
			Right: scene.Position3D{X: 2, Y: 1, Z: 1},
		},
		Receiver: scene.ReceiverConfig{X: 1.5, Y: 3, Z: 1},
	}
	catalog := materials.GetCatalog()
	
	errors := ValidateScene(config, catalog)
	
	if errors["room.widthM"] == "" {
		t.Error("Expected error for negative width")
	}
}

func TestValidateScene_SourceOutOfBounds(t *testing.T) {
	config := scene.SceneConfig{
		Room: scene.RoomConfig{
			WidthM:  4,
			LengthM: 6,
			HeightM: 3,
		},
		Surfaces: scene.SurfaceConfig{
			FrontWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RearWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			LeftWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RightWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			Ceiling:    scene.SurfaceMaterial{MaterialID: "painted_concrete"},
			Floor:      scene.SurfaceMaterial{MaterialID: "carpet"},
		},
		Sources: scene.SourcesConfig{
			Left:  scene.Position3D{X: 10, Y: 1, Z: 1}, // Out of bounds
			Right: scene.Position3D{X: 2, Y: 1, Z: 1},
		},
		Receiver: scene.ReceiverConfig{X: 1.5, Y: 3, Z: 1},
	}
	catalog := materials.GetCatalog()
	
	errors := ValidateScene(config, catalog)
	
	if errors["sources.left"] == "" {
		t.Error("Expected error for out of bounds source")
	}
}

func TestValidateScene_ReceiverOutOfBounds(t *testing.T) {
	config := scene.SceneConfig{
		Room: scene.RoomConfig{
			WidthM:  4,
			LengthM: 6,
			HeightM: 3,
		},
		Surfaces: scene.SurfaceConfig{
			FrontWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RearWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			LeftWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RightWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			Ceiling:    scene.SurfaceMaterial{MaterialID: "painted_concrete"},
			Floor:      scene.SurfaceMaterial{MaterialID: "carpet"},
		},
		Sources: scene.SourcesConfig{
			Left:  scene.Position3D{X: 1, Y: 1, Z: 1},
			Right: scene.Position3D{X: 2, Y: 1, Z: 1},
		},
		Receiver: scene.ReceiverConfig{X: 10, Y: 3, Z: 1}, // Out of bounds
	}
	catalog := materials.GetCatalog()
	
	errors := ValidateScene(config, catalog)
	
	if errors["receiver"] == "" {
		t.Error("Expected error for out of bounds receiver")
	}
}

func TestValidateScene_InvalidMaterial(t *testing.T) {
	config := scene.SceneConfig{
		Room: scene.RoomConfig{
			WidthM:  4,
			LengthM: 6,
			HeightM: 3,
		},
		Surfaces: scene.SurfaceConfig{
			FrontWall:  scene.SurfaceMaterial{MaterialID: "invalid_material"},
			RearWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			LeftWall:   scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			RightWall:  scene.SurfaceMaterial{MaterialID: "gypsum_board"},
			Ceiling:    scene.SurfaceMaterial{MaterialID: "painted_concrete"},
			Floor:      scene.SurfaceMaterial{MaterialID: "carpet"},
		},
		Sources: scene.SourcesConfig{
			Left:  scene.Position3D{X: 1, Y: 1, Z: 1},
			Right: scene.Position3D{X: 2, Y: 1, Z: 1},
		},
		Receiver: scene.ReceiverConfig{X: 1.5, Y: 3, Z: 1},
	}
	catalog := materials.GetCatalog()
	
	errors := ValidateScene(config, catalog)
	
	if errors["surfaces.frontWall"] == "" {
		t.Error("Expected error for invalid material")
	}
}
