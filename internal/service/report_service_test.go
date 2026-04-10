package service

import (
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestProjectService_CreateDemoProject(t *testing.T) {
	service := &ProjectService{}

	input := domain.CreateProjectInput{
		Name:            "Demo Studio Project",
		SpaceType:       domain.SpaceTypeStudio,
		GoalType:        domain.GoalTypeMixing,
		Units:           domain.UnitsMetric,
		TemperatureC:    20.0,
		HumidityPercent: 50.0,
	}

	if err := service.validateCreateInput(input); err != nil {
		t.Errorf("Demo project input should be valid: %v", err)
	}

	if input.Name != "Demo Studio Project" {
		t.Errorf("Expected name 'Demo Studio Project', got '%s'", input.Name)
	}

	if input.SpaceType != domain.SpaceTypeStudio {
		t.Errorf("Expected space type '%s', got '%s'", domain.SpaceTypeStudio, input.SpaceType)
	}

	if input.GoalType != domain.GoalTypeMixing {
		t.Errorf("Expected goal type '%s', got '%s'", domain.GoalTypeMixing, input.GoalType)
	}
}
