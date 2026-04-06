package domain

import (
	"testing"
	"time"
)

func TestNewProject(t *testing.T) {
	input := CreateProjectInput{
		Name:            "Test Project",
		SpaceType:       SpaceTypeStudio,
		GoalType:        GoalTypeMixing,
		Units:           UnitsMetric,
		TemperatureC:    20.0,
		HumidityPercent: 50.0,
		SpeedOfSound:    343.0,
	}

	project := NewProject(input)

	if project.ID == "" {
		t.Error("Expected project ID to be generated")
	}

	if project.Name != input.Name {
		t.Errorf("Expected name %s, got %s", input.Name, project.Name)
	}

	if project.SpaceType != input.SpaceType {
		t.Errorf("Expected space type %s, got %s", input.SpaceType, project.SpaceType)
	}

	if project.GoalType != input.GoalType {
		t.Errorf("Expected goal type %s, got %s", input.GoalType, project.GoalType)
	}

	if project.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if project.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if project.CreatedAt != project.UpdatedAt {
		t.Error("Expected CreatedAt and UpdatedAt to be equal for new project")
	}
}

func TestProject_Update(t *testing.T) {
	project := NewProject(CreateProjectInput{
		Name:      "Original",
		SpaceType: SpaceTypeStudio,
		GoalType:  GoalTypeMixing,
		Units:     UnitsMetric,
	})

	originalUpdatedAt := project.UpdatedAt
	
	// Wait a bit to ensure different timestamp
	time.Sleep(1 * time.Millisecond)

	newName := "Updated"
	input := UpdateProjectInput{
		Name: &newName,
	}

	project.Update(input)

	if project.Name != newName {
		t.Errorf("Expected name %s, got %s", newName, project.Name)
	}

	if !project.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be after original UpdatedAt")
	}
}

func TestSpaceType_DisplayName(t *testing.T) {
	tests := map[SpaceType]string{
		SpaceTypeStudio: "Studio / Control Room",
		SpaceTypeHifi:   "Hi-fi Listening Room",
		SpaceTypePublic: "Public Space",
	}

	for spaceType, expected := range tests {
		if spaceType.DisplayName() != expected {
			t.Errorf("Expected %s to display as %s, got %s", spaceType, expected, spaceType.DisplayName())
		}
	}
}

func TestSpaceType_Description(t *testing.T) {
	tests := map[SpaceType]string{
		SpaceTypeStudio: "Critical listening environment with symmetry and early reflection control",
		SpaceTypeHifi:   "Optimized listening room for imaging and spaciousness",
		SpaceTypePublic: "Audience-focused space with coverage and clarity requirements",
	}

	for spaceType, expected := range tests {
		if spaceType.Description() != expected {
			t.Errorf("Expected %s description to be %s, got %s", spaceType, expected, spaceType.Description())
		}
	}
}

func TestGoalType_DisplayName(t *testing.T) {
	tests := map[GoalType]string{
		GoalTypeMixing:     "Mixing",
		GoalTypeMastering:  "Mastering",
		GoalTypeHifiStereo: "Hi-fi Stereo",
		GoalTypeSpeech:     "Speech",
		GoalTypeMusic:      "Music",
		GoalTypeMultipurpose: "Multipurpose",
	}

	for goalType, expected := range tests {
		if goalType.DisplayName() != expected {
			t.Errorf("Expected %s to display as %s, got %s", goalType, expected, goalType.DisplayName())
		}
	}
}

func TestValidSpaceTypes(t *testing.T) {
	types := ValidSpaceTypes()
	
	expected := []SpaceType{SpaceTypeStudio, SpaceTypeHifi, SpaceTypePublic}
	
	if len(types) != len(expected) {
		t.Errorf("Expected %d space types, got %d", len(expected), len(types))
	}
	
	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("Expected space type %d to be %s, got %s", i, expectedType, types[i])
		}
	}
}

func TestValidGoalTypes(t *testing.T) {
	types := ValidGoalTypes()
	
	expected := []GoalType{
		GoalTypeMixing,
		GoalTypeMastering,
		GoalTypeHifiStereo,
		GoalTypeSpeech,
		GoalTypeMusic,
		GoalTypeMultipurpose,
	}
	
	if len(types) != len(expected) {
		t.Errorf("Expected %d goal types, got %d", len(expected), len(types))
	}
	
	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("Expected goal type %d to be %s, got %s", i, expectedType, types[i])
		}
	}
}

func TestValidUnits(t *testing.T) {
	types := ValidUnits()
	
	expected := []Units{UnitsMetric, UnitsImperial}
	
	if len(types) != len(expected) {
		t.Errorf("Expected %d units, got %d", len(expected), len(types))
	}
	
	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("Expected unit %d to be %s, got %s", i, expectedType, types[i])
		}
	}
}
