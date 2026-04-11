package validation

import (
	"acoustics-calculator/internal/materials"
	"acoustics-calculator/internal/scene"
)

func ValidateScene(config scene.SceneConfig, catalog map[string]materials.Material) map[string]string {
	errors := make(map[string]string)

	// Validate room dimensions
	if config.Room.WidthM <= 0 {
		errors["room.widthM"] = "must be greater than zero"
	}
	if config.Room.LengthM <= 0 {
		errors["room.lengthM"] = "must be greater than zero"
	}
	if config.Room.HeightM <= 0 {
		errors["room.heightM"] = "must be greater than zero"
	}

	// Validate sources are inside room bounds
	if !isInsideRoom(config.Sources.Left.X, config.Sources.Left.Y, config.Sources.Left.Z, config.Room) {
		errors["sources.left"] = "must be inside room bounds"
	}
	if !isInsideRoom(config.Sources.Right.X, config.Sources.Right.Y, config.Sources.Right.Z, config.Room) {
		errors["sources.right"] = "must be inside room bounds"
	}

	// Validate receiver is inside room bounds
	if !isInsideRoom(config.Receiver.X, config.Receiver.Y, config.Receiver.Z, config.Room) {
		errors["receiver"] = "must be inside room bounds"
	}

	// Validate materials exist
	validateMaterial(config.Surfaces.FrontWall.MaterialID, "surfaces.frontWall", catalog, errors)
	validateMaterial(config.Surfaces.RearWall.MaterialID, "surfaces.rearWall", catalog, errors)
	validateMaterial(config.Surfaces.LeftWall.MaterialID, "surfaces.leftWall", catalog, errors)
	validateMaterial(config.Surfaces.RightWall.MaterialID, "surfaces.rightWall", catalog, errors)
	validateMaterial(config.Surfaces.Ceiling.MaterialID, "surfaces.ceiling", catalog, errors)
	validateMaterial(config.Surfaces.Floor.MaterialID, "surfaces.floor", catalog, errors)

	return errors
}

func isInsideRoom(x, y, z float64, room scene.RoomConfig) bool {
	return x >= 0 && x <= room.WidthM &&
		y >= 0 && y <= room.LengthM &&
		z >= 0 && z <= room.HeightM
}

func validateMaterial(materialID, field string, catalog map[string]materials.Material, errors map[string]string) {
	if _, ok := catalog[materialID]; !ok {
		errors[field] = "material not found in catalog"
	}
}
