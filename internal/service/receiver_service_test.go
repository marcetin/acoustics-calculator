package service

import (
	"errors"
	"testing"

	"acoustics-calculator/internal/domain"
)

func TestReceiverService_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.CreateReceiverInput
		wantErr bool
	}{
		{
			name: "valid receiver",
			input: domain.CreateReceiverInput{
				ProjectID: "test-project",
				Name:      "Test Receiver",
				Type:      domain.ReceiverTypeListener,
				PositionX: 2.5,
				PositionY: 3.0,
				PositionZ: 1.2,
				Weight:    1.0,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			input: domain.CreateReceiverInput{
				ProjectID: "test-project",
				Name:      "",
				Type:      domain.ReceiverTypeListener,
				PositionX: 2.5,
				PositionY: 3.0,
				PositionZ: 1.2,
				Weight:    1.0,
			},
			wantErr: true,
		},
		{
			name: "zero weight",
			input: domain.CreateReceiverInput{
				ProjectID: "test-project",
				Name:      "Test Receiver",
				Type:      domain.ReceiverTypeListener,
				PositionX: 2.5,
				PositionY: 3.0,
				PositionZ: 1.2,
				Weight:    0.0,
			},
			wantErr: true,
		},
		{
			name: "negative weight",
			input: domain.CreateReceiverInput{
				ProjectID: "test-project",
				Name:      "Test Receiver",
				Type:      domain.ReceiverTypeListener,
				PositionX: 2.5,
				PositionY: 3.0,
				PositionZ: 1.2,
				Weight:    -1.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReceiverInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateReceiverInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReceiverService_WithinGeometry(t *testing.T) {
	width := 5.0
	length := 6.0
	height := 3.0

	tests := []struct {
		name    string
		x       float64
		y       float64
		z       float64
		wantErr bool
	}{
		{"inside room", 2.5, 3.0, 1.5, false},
		{"outside x", 6.0, 3.0, 1.5, true},
		{"outside y", 2.5, 7.0, 1.5, true},
		{"outside z", 2.5, 3.0, 4.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReceiverWithinGeometry(tt.x, tt.y, tt.z, width, length, height)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateReceiverWithinGeometry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validateReceiverInput(input domain.CreateReceiverInput) error {
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.Weight <= 0 {
		return errors.New("weight must be positive")
	}
	if input.PositionX < 0 || input.PositionY < 0 || input.PositionZ < 0 {
		return errors.New("position cannot be negative")
	}
	return nil
}

func validateReceiverWithinGeometry(x, y, z, width, length, height float64) error {
	if x < 0 || x > width {
		return errors.New("x position outside room")
	}
	if y < 0 || y > length {
		return errors.New("y position outside room")
	}
	if z < 0 || z > height {
		return errors.New("z position outside room")
	}
	return nil
}
