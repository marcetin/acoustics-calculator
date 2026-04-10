package acoustics

import (
	"testing"
)

func TestModalSolver_CalculateModes(t *testing.T) {
	tests := []struct {
		name     string
		width    float64
		length   float64
		height   float64
		maxFreq  float64
		c        float64
		wantMin  int
		wantKind bool
	}{
		{
			name:     "simple room",
			width:    5.0,
			length:   6.0,
			height:   3.0,
			maxFreq:  300,
			c:        343,
			wantMin:  1,
			wantKind: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solver := NewModalSolver(tt.width, tt.length, tt.height, tt.c)
			modes := solver.CalculateModes(tt.maxFreq)

			if len(modes) < tt.wantMin {
				t.Errorf("CalculateModes() returned %d modes, want at least %d", len(modes), tt.wantMin)
			}

			if tt.wantKind && len(modes) > 0 {
				validKind := modes[0].Kind == "AXIAL" || modes[0].Kind == "TANGENTIAL" || modes[0].Kind == "OBLIQUE"
				if !validKind {
					t.Errorf("CalculateModes() returned invalid kind: %s", modes[0].Kind)
				}
			}
		})
	}
}

func TestModalSolver_ModeClassification(t *testing.T) {
	solver := NewModalSolver(5.0, 6.0, 3.0, 343)
	modes := solver.CalculateModes(100)

	axialCount := 0
	tangentialCount := 0
	obliqueCount := 0

	for _, mode := range modes {
		switch mode.Kind {
		case "AXIAL":
			axialCount++
		case "TANGENTIAL":
			tangentialCount++
		case "OBLIQUE":
			obliqueCount++
		}
	}

	if axialCount == 0 && tangentialCount == 0 && obliqueCount == 0 {
		t.Error("No modes classified correctly")
	}
}

func TestModalSolver_SortedByFrequency(t *testing.T) {
	solver := NewModalSolver(5.0, 6.0, 3.0, 343)
	modes := solver.CalculateModes(300)

	for i := 1; i < len(modes); i++ {
		if modes[i].Freq < modes[i-1].Freq {
			t.Errorf("Modes not sorted: mode %d freq %.2f < mode %d freq %.2f", i, modes[i].Freq, i-1, modes[i-1].Freq)
		}
	}
}

func TestModalSolver_CalculateReceiverPressure(t *testing.T) {
	solver := NewModalSolver(5.0, 6.0, 3.0, 343)
	modes := solver.CalculateModes(300)

	pressure := solver.CalculateReceiverPressure(2.5, 3.0, 1.5, modes)

	if pressure < 0 {
		t.Errorf("CalculateReceiverPressure() returned negative pressure: %f", pressure)
	}
}

func TestCalculateSpeedOfSound(t *testing.T) {
	tests := []struct {
		name     string
		tempC    float64
		expected float64
	}{
		{"20C", 20.0, 331.3 + 0.606*20.0},
		{"0C", 0.0, 331.3},
		{"extreme hot", 150.0, SpeedOfSoundDefault},
		{"extreme cold", -100.0, SpeedOfSoundDefault},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSpeedOfSound(tt.tempC)
			if got != tt.expected && tt.tempC >= -50 && tt.tempC <= 100 {
				t.Errorf("CalculateSpeedOfSound() = %v, want %v", got, tt.expected)
			}
			if tt.tempC < -50 || tt.tempC > 100 {
				if got != SpeedOfSoundDefault {
					t.Errorf("CalculateSpeedOfSound() should use default for extreme temps, got %v", got)
				}
			}
		})
	}
}
