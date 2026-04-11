package scene

import "testing"

func TestGetExample(t *testing.T) {
	config := GetExample()
	
	if config.Room.WidthM <= 0 {
		t.Error("Expected positive room width")
	}
	if config.Room.LengthM <= 0 {
		t.Error("Expected positive room length")
	}
	if config.Room.HeightM <= 0 {
		t.Error("Expected positive room height")
	}
	
	if config.Surfaces.FrontWall.MaterialID == "" {
		t.Error("Expected front wall material ID")
	}
	if config.Surfaces.RearWall.MaterialID == "" {
		t.Error("Expected rear wall material ID")
	}
	if config.Surfaces.LeftWall.MaterialID == "" {
		t.Error("Expected left wall material ID")
	}
	if config.Surfaces.RightWall.MaterialID == "" {
		t.Error("Expected right wall material ID")
	}
	if config.Surfaces.Ceiling.MaterialID == "" {
		t.Error("Expected ceiling material ID")
	}
	if config.Surfaces.Floor.MaterialID == "" {
		t.Error("Expected floor material ID")
	}
	
	if config.Sources.Left.X < 0 || config.Sources.Left.X > config.Room.WidthM {
		t.Error("Left speaker X out of bounds")
	}
	if config.Sources.Right.X < 0 || config.Sources.Right.X > config.Room.WidthM {
		t.Error("Right speaker X out of bounds")
	}
	
	if config.Receiver.X < 0 || config.Receiver.X > config.Room.WidthM {
		t.Error("Receiver X out of bounds")
	}
}
