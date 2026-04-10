package service

import (
	"context"
	"fmt"
	"strings"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type ProjectService struct {
	projectRepo *repo.ProjectRepository
}

func NewProjectService(projectRepo *repo.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, input domain.CreateProjectInput) (*domain.Project, error) {
	if err := s.validateCreateInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	normalizedInput := s.setDefaults(input)

	project := domain.NewProject(normalizedInput)

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	if id == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) ListProjects(ctx context.Context) ([]*domain.Project, error) {
	projects, err := s.projectRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, input domain.UpdateProjectInput) (*domain.Project, error) {
	if id == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	if err := s.validateUpdateInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	project.Update(input)

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) validateCreateInput(input domain.CreateProjectInput) error {
	var errors []string

	if strings.TrimSpace(input.Name) == "" {
		errors = append(errors, "project name is required")
	}

	if !s.isValidSpaceType(input.SpaceType) {
		errors = append(errors, "invalid space type")
	}

	if !s.isValidGoalType(input.GoalType) {
		errors = append(errors, "invalid goal type")
	}

	if !s.isValidUnits(input.Units) {
		errors = append(errors, "invalid units")
	}

	if input.TemperatureC < -273.15 || input.TemperatureC > 100 {
		errors = append(errors, "temperature must be between -273.15°C and 100°C")
	}

	if input.HumidityPercent < 0 || input.HumidityPercent > 100 {
		errors = append(errors, "humidity must be between 0% and 100%")
	}

	if input.SpeedOfSound < 0 || input.SpeedOfSound > 1000 {
		errors = append(errors, "speed of sound must be between 0 and 1000 m/s")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

func (s *ProjectService) validateUpdateInput(input domain.UpdateProjectInput) error {
	var errors []string

	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		errors = append(errors, "project name cannot be empty")
	}

	if input.SpaceType != nil && !s.isValidSpaceType(*input.SpaceType) {
		errors = append(errors, "invalid space type")
	}

	if input.GoalType != nil && !s.isValidGoalType(*input.GoalType) {
		errors = append(errors, "invalid goal type")
	}

	if input.Units != nil && !s.isValidUnits(*input.Units) {
		errors = append(errors, "invalid units")
	}

	if input.TemperatureC != nil && (*input.TemperatureC < -273.15 || *input.TemperatureC > 100) {
		errors = append(errors, "temperature must be between -273.15°C and 100°C")
	}

	if input.HumidityPercent != nil && (*input.HumidityPercent < 0 || *input.HumidityPercent > 100) {
		errors = append(errors, "humidity must be between 0% and 100%")
	}

	if input.SpeedOfSound != nil && (*input.SpeedOfSound < 0 || *input.SpeedOfSound > 1000) {
		errors = append(errors, "speed of sound must be between 0 and 1000 m/s")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

func (s *ProjectService) setDefaults(input domain.CreateProjectInput) domain.CreateProjectInput {
	if input.Units == "" {
		input.Units = domain.UnitsMetric
	}

	if input.TemperatureC == 0 {
		input.TemperatureC = 20.0
	}

	if input.HumidityPercent == 0 {
		input.HumidityPercent = 50.0
	}

	if input.SpeedOfSound == 0 {
		input.SpeedOfSound = 343.0
	}

	return input
}

func (s *ProjectService) isValidSpaceType(spaceType domain.SpaceType) bool {
	for _, validType := range domain.ValidSpaceTypes() {
		if spaceType == validType {
			return true
		}
	}
	return false
}

func (s *ProjectService) isValidGoalType(goalType domain.GoalType) bool {
	for _, validType := range domain.ValidGoalTypes() {
		if goalType == validType {
			return true
		}
	}
	return false
}

func (s *ProjectService) isValidUnits(units domain.Units) bool {
	for _, validUnits := range domain.ValidUnits() {
		if units == validUnits {
			return true
		}
	}
	return false
}

func (s *ProjectService) CreateDemoProject(ctx context.Context) (*domain.Project, error) {
	input := domain.CreateProjectInput{
		Name:            "Demo Studio Project",
		SpaceType:       domain.SpaceTypeStudio,
		GoalType:        domain.GoalTypeMixing,
		Units:           domain.UnitsMetric,
		TemperatureC:    20.0,
		HumidityPercent: 50.0,
	}

	return s.CreateProject(ctx, input)
}
