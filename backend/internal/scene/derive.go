package scene

import "acoustics-calculator/internal/materials"

func DeriveScene(config SceneConfig, catalog map[string]materials.Material) (DerivedScene, []string) {
	warnings := []string{}

	// Calculate room bounds
	bounds := RoomBounds{
		MinX: 0,
		MaxX: config.Room.WidthM,
		MinY: 0,
		MaxY: config.Room.LengthM,
		MinZ: 0,
		MaxZ: config.Room.HeightM,
	}

	// Derive surfaces
	surfaces := []DerivedSurface{
		{
			ID:         "front",
			Label:      "Front Wall",
			MaterialID: config.Surfaces.FrontWall.MaterialID,
			ColorHex:   getMaterialColor(config.Surfaces.FrontWall.MaterialID, catalog),
		},
		{
			ID:         "rear",
			Label:      "Rear Wall",
			MaterialID: config.Surfaces.RearWall.MaterialID,
			ColorHex:   getMaterialColor(config.Surfaces.RearWall.MaterialID, catalog),
		},
		{
			ID:         "left",
			Label:      "Left Wall",
			MaterialID: config.Surfaces.LeftWall.MaterialID,
			ColorHex:   getMaterialColor(config.Surfaces.LeftWall.MaterialID, catalog),
		},
		{
			ID:         "right",
			Label:      "Right Wall",
			MaterialID: config.Surfaces.RightWall.MaterialID,
			ColorHex:   getMaterialColor(config.Surfaces.RightWall.MaterialID, catalog),
		},
		{
			ID:         "ceiling",
			Label:      "Ceiling",
			MaterialID: config.Surfaces.Ceiling.MaterialID,
			ColorHex:   getMaterialColor(config.Surfaces.Ceiling.MaterialID, catalog),
		},
		{
			ID:         "floor",
			Label:      "Floor",
			MaterialID: config.Surfaces.Floor.MaterialID,
			ColorHex:   getMaterialColor(config.Surfaces.Floor.MaterialID, catalog),
		},
	}

	derived := DerivedScene{
		Surfaces:     surfaces,
		RoomBounds:   bounds,
		LeftSpeaker:  config.Sources.Left,
		RightSpeaker: config.Sources.Right,
		Receiver:     Position3D{X: config.Receiver.X, Y: config.Receiver.Y, Z: config.Receiver.Z},
	}

	return derived, warnings
}

func getMaterialColor(materialID string, catalog map[string]materials.Material) string {
	if mat, ok := catalog[materialID]; ok {
		return mat.ColorHex
	}
	return "#CCCCCC"
}
