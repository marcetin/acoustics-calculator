package service

import (
	"context"
	"fmt"

	"acoustics-calculator/internal/domain"
	"acoustics-calculator/internal/repo"
)

type ReceiverService struct {
	receiverRepo *repo.ReceiverRepository
	geometryRepo *repo.GeometryRepository
}

func NewReceiverService(receiverRepo *repo.ReceiverRepository, geometryRepo *repo.GeometryRepository) *ReceiverService {
	return &ReceiverService{
		receiverRepo: receiverRepo,
		geometryRepo: geometryRepo,
	}
}

func (s *ReceiverService) ListByProject(ctx context.Context, projectID string) ([]*domain.Receiver, error) {
	return s.receiverRepo.ListByProjectID(ctx, projectID)
}

func (s *ReceiverService) GetByID(ctx context.Context, id string) (*domain.Receiver, error) {
	return s.receiverRepo.GetByID(ctx, id)
}

func (s *ReceiverService) Create(ctx context.Context, input domain.CreateReceiverInput) (*domain.Receiver, error) {
	if err := s.validateReceiverInput(ctx, input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	receiver := domain.NewReceiver(input)

	if err := s.receiverRepo.Create(ctx, receiver); err != nil {
		return nil, fmt.Errorf("failed to create receiver: %w", err)
	}

	return receiver, nil
}

func (s *ReceiverService) Update(ctx context.Context, id string, input domain.UpdateReceiverInput) (*domain.Receiver, error) {
	receiver, err := s.receiverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get receiver: %w", err)
	}

	if err := s.validateReceiverUpdate(ctx, receiver.ProjectID, input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	receiver.Update(input)

	if err := s.receiverRepo.Update(ctx, receiver); err != nil {
		return nil, fmt.Errorf("failed to update receiver: %w", err)
	}

	return receiver, nil
}

func (s *ReceiverService) Delete(ctx context.Context, id string) error {
	return s.receiverRepo.Delete(ctx, id)
}

func (s *ReceiverService) validateReceiverInput(ctx context.Context, input domain.CreateReceiverInput) error {
	var errors []string

	if input.ProjectID == "" {
		errors = append(errors, "project ID is required")
	}

	if input.Name == "" {
		errors = append(errors, "receiver name is required")
	}

	if input.Weight <= 0 {
		errors = append(errors, "weight must be greater than 0")
	}

	if input.PositionZ < 0 {
		errors = append(errors, "Z position cannot be negative")
	}

	if err := s.validateWithinGeometry(ctx, input.ProjectID, input.PositionX, input.PositionY, input.PositionZ); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}

func (s *ReceiverService) validateReceiverUpdate(ctx context.Context, projectID string, input domain.UpdateReceiverInput) error {
	var errors []string

	if input.Weight != nil && *input.Weight <= 0 {
		errors = append(errors, "weight must be greater than 0")
	}

	if input.PositionX != nil && input.PositionY != nil && input.PositionZ != nil {
		if *input.PositionZ < 0 {
			errors = append(errors, "Z position cannot be negative")
		}

		if err := s.validateWithinGeometry(ctx, projectID, *input.PositionX, *input.PositionY, *input.PositionZ); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}

	return nil
}

func (s *ReceiverService) validateWithinGeometry(ctx context.Context, projectID string, x, y, z float64) error {
	geometry, err := s.geometryRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to get geometry: %w", err)
	}

	if geometry == nil {
		return fmt.Errorf("room geometry must be defined before adding receivers")
	}

	if !geometry.IsShoebox() {
		return fmt.Errorf("position validation not yet implemented for polyhedral geometry")
	}

	if x < 0 || x > geometry.Width {
		return fmt.Errorf("X position must be between 0 and room width")
	}

	if y < 0 || y > geometry.Length {
		return fmt.Errorf("Y position must be between 0 and room length")
	}

	if z < 0 || z > geometry.Height {
		return fmt.Errorf("Z position must be between 0 and room height")
	}

	return nil
}
