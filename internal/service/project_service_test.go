package service

import (
	"context"
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestProjectService_ValidateCreateInput(t *testing.T) {
	service := &ProjectService{}

	tests := []struct {
		name    string
		input   domain.CreateProjectInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: domain.CreateProjectInput{
				Name:            "Test Project",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsMetric,
				TemperatureC:    20.0,
				HumidityPercent: 50.0,
				SpeedOfSound:    343.0,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: domain.CreateProjectInput{
				Name:            "",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsMetric,
			},
			wantErr: true,
		},
		{
			name: "invalid space type",
			input: domain.CreateProjectInput{
				Name:            "Test",
				SpaceType:       "INVALID",
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsMetric,
			},
			wantErr: true,
		},
		{
			name: "temperature too low",
			input: domain.CreateProjectInput{
				Name:            "Test",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsMetric,
				TemperatureC:    -300,
			},
			wantErr: true,
		},
		{
			name: "humidity too high",
			input: domain.CreateProjectInput{
				Name:            "Test",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsMetric,
				HumidityPercent: 150,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCreateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectService_SetDefaults(t *testing.T) {
	service := &ProjectService{}

	tests := []struct {
		name  string
		input domain.CreateProjectInput
		want  domain.CreateProjectInput
	}{
		{
			name: "all defaults",
			input: domain.CreateProjectInput{
				Name:      "Test",
				SpaceType: domain.SpaceTypeStudio,
				GoalType:  domain.GoalTypeMixing,
			},
			want: domain.CreateProjectInput{
				Name:            "Test",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsMetric,
				TemperatureC:    20.0,
				HumidityPercent: 50.0,
				SpeedOfSound:    343.0,
			},
		},
		{
			name: "some values set",
			input: domain.CreateProjectInput{
				Name:            "Test",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsImperial,
				TemperatureC:    25.0,
			},
			want: domain.CreateProjectInput{
				Name:            "Test",
				SpaceType:       domain.SpaceTypeStudio,
				GoalType:        domain.GoalTypeMixing,
				Units:           domain.UnitsImperial,
				TemperatureC:    25.0,
				HumidityPercent: 50.0,
				SpeedOfSound:    343.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.setDefaults(tt.input)
			if got.Units != tt.want.Units ||
				got.TemperatureC != tt.want.TemperatureC ||
				got.HumidityPercent != tt.want.HumidityPercent ||
				got.SpeedOfSound != tt.want.SpeedOfSound {
				t.Errorf("setDefaults() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestProjectService_GetProject_InvalidID(t *testing.T) {
	// This would require a mock repository for full testing
	// For now, just test the validation logic
	service := &ProjectService{}

	// Test with empty context and empty ID
	_, err := service.GetProject(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty project ID")
	}
}
