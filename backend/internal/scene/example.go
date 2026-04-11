package scene

func GetExample() SceneConfig {
	return SceneConfig{
		Room: RoomConfig{
			WidthM:  4.2,
			LengthM: 6.1,
			HeightM: 2.8,
		},
		Surfaces: SurfaceConfig{
			FrontWall: SurfaceMaterial{MaterialID: "gypsum_board"},
			RearWall:  SurfaceMaterial{MaterialID: "wood_panel"},
			LeftWall:  SurfaceMaterial{MaterialID: "gypsum_board"},
			RightWall: SurfaceMaterial{MaterialID: "gypsum_board"},
			Ceiling:   SurfaceMaterial{MaterialID: "painted_concrete"},
			Floor:     SurfaceMaterial{MaterialID: "carpet"},
		},
		Sources: SourcesConfig{
			Left:  Position3D{X: 1.3, Y: 1.1, Z: 1.2},
			Right: Position3D{X: 2.9, Y: 1.1, Z: 1.2},
		},
		Receiver: ReceiverConfig{
			X: 2.1,
			Y: 3.3,
			Z: 1.2,
		},
		Viewer: ViewerConfig{
			ShowGrid:   true,
			ShowLabels: true,
			ShowAxes:   true,
		},
	}
}
