package materials

type Material struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	ColorHex string `json:"colorHex"`
}

func GetCatalog() map[string]Material {
	return map[string]Material{
		"gypsum_board": {
			ID:       "gypsum_board",
			Label:    "Gypsum Board",
			ColorHex: "#E8E8E8",
		},
		"painted_concrete": {
			ID:       "painted_concrete",
			Label:    "Painted Concrete",
			ColorHex: "#B0B0B0",
		},
		"wood_panel": {
			ID:       "wood_panel",
			Label:    "Wood Panel",
			ColorHex: "#8B6914",
		},
		"carpet": {
			ID:       "carpet",
			Label:    "Carpet",
			ColorHex: "#4A4A4A",
		},
		"glass": {
			ID:       "glass",
			Label:    "Glass",
			ColorHex: "#A8D8EA",
		},
		"curtain_heavy": {
			ID:       "curtain_heavy",
			Label:    "Heavy Curtain",
			ColorHex: "#5C4033",
		},
	}
}
