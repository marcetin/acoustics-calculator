package acoustics

import (
	"math"
	"sort"

	"acoustics-calculator/internal/domain"
)

const SpeedOfSoundDefault = 343.0 // m/s at 20°C

func CalculateSpeedOfSound(tempC float64) float64 {
	if tempC < -50 || tempC > 100 {
		return SpeedOfSoundDefault
	}
	return 331.3 + 0.606*tempC
}

type ModalSolver struct {
	width  float64
	length float64
	height float64
	c      float64
}

func NewModalSolver(width, length, height, c float64) *ModalSolver {
	return &ModalSolver{
		width:  width,
		length: length,
		height: height,
		c:      c,
	}
}

func (ms *ModalSolver) CalculateModes(maxFreq float64) []domain.Mode {
	var modes []domain.Mode

	lMax := int(maxFreq * 2 * ms.width / ms.c)
	mMax := int(maxFreq * 2 * ms.length / ms.c)
	nMax := int(maxFreq * 2 * ms.height / ms.c)

	for l := 0; l <= lMax; l++ {
		for m := 0; m <= mMax; m++ {
			for n := 0; n <= nMax; n++ {
				if l == 0 && m == 0 && n == 0 {
					continue
				}

				freq := ms.calculateFrequency(l, m, n)
				if freq > maxFreq {
					continue
				}

				kind := ms.classifyMode(l, m, n)

				modes = append(modes, domain.Mode{
					L:    l,
					M:    m,
					N:    n,
					Freq: freq,
					Kind: kind,
				})
			}
		}
	}

	sort.Slice(modes, func(i, j int) bool {
		return modes[i].Freq < modes[j].Freq
	})

	return modes
}

func (ms *ModalSolver) calculateFrequency(l, m, n int) float64 {
	return (ms.c / 2) * math.Sqrt(
		math.Pow(float64(l)/ms.width, 2)+
			math.Pow(float64(m)/ms.length, 2)+
			math.Pow(float64(n)/ms.height, 2),
	)
}

func (ms *ModalSolver) classifyMode(l, m, n int) string {
	nonZeroCount := 0
	if l > 0 {
		nonZeroCount++
	}
	if m > 0 {
		nonZeroCount++
	}
	if n > 0 {
		nonZeroCount++
	}

	switch nonZeroCount {
	case 1:
		return "AXIAL"
	case 2:
		return "TANGENTIAL"
	default:
		return "OBLIQUE"
	}
}

func (ms *ModalSolver) CalculateReceiverPressure(x, y, z float64, modes []domain.Mode) float64 {
	if len(modes) == 0 {
		return 0
	}

	totalPressure := 0.0
	for _, mode := range modes {
		if mode.Freq <= 0 {
			continue
		}

		cosL := math.Cos(float64(mode.L) * math.Pi * x / ms.width)
		cosM := math.Cos(float64(mode.M) * math.Pi * y / ms.length)
		cosN := math.Cos(float64(mode.N) * math.Pi * z / ms.height)

		pressure := cosL * cosM * cosN / mode.Freq
		totalPressure += math.Abs(pressure)
	}

	return totalPressure
}
