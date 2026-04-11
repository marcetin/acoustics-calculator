package scene

type SceneConfig struct {
	Room      RoomConfig      `json:"room"`
	Surfaces  SurfaceConfig   `json:"surfaces"`
	Sources   SourcesConfig   `json:"sources"`
	Receiver  ReceiverConfig  `json:"receiver"`
	Viewer    ViewerConfig    `json:"viewer"`
}

type RoomConfig struct {
	WidthM  float64 `json:"widthM"`
	LengthM float64 `json:"lengthM"`
	HeightM float64 `json:"heightM"`
}

type SurfaceConfig struct {
	FrontWall  SurfaceMaterial `json:"frontWall"`
	RearWall   SurfaceMaterial `json:"rearWall"`
	LeftWall   SurfaceMaterial `json:"leftWall"`
	RightWall  SurfaceMaterial `json:"rightWall"`
	Ceiling    SurfaceMaterial `json:"ceiling"`
	Floor      SurfaceMaterial `json:"floor"`
}

type SurfaceMaterial struct {
	MaterialID string `json:"materialId"`
}

type SourcesConfig struct {
	Left  Position3D `json:"left"`
	Right Position3D `json:"right"`
}

type Position3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ReceiverConfig struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ViewerConfig struct {
	ShowGrid   bool `json:"showGrid"`
	ShowLabels bool `json:"showLabels"`
	ShowAxes   bool `json:"showAxes"`
}

type DerivedScene struct {
	Surfaces      []DerivedSurface `json:"surfaces"`
	RoomBounds    RoomBounds        `json:"roomBounds"`
	LeftSpeaker   Position3D        `json:"leftSpeaker"`
	RightSpeaker  Position3D        `json:"rightSpeaker"`
	Receiver      Position3D        `json:"receiver"`
}

type DerivedSurface struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	MaterialID string  `json:"materialId"`
	ColorHex   string  `json:"colorHex"`
}

type RoomBounds struct {
	MinX float64 `json:"minX"`
	MaxX float64 `json:"maxX"`
	MinY float64 `json:"minY"`
	MaxY float64 `json:"maxY"`
	MinZ float64 `json:"minZ"`
	MaxZ float64 `json:"maxZ"`
}
