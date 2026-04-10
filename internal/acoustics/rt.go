package acoustics

import (
	"fmt"
	"math"

	"acoustics-calculator/internal/domain"
)

type RTEstimator struct {
	volume float64
	c      float64
}

func NewRTEstimator(volume, c float64) *RTEstimator {
	return &RTEstimator{
		volume: volume,
		c:      c,
	}
}

var OctaveBands = []float64{63, 125, 250, 500, 1000, 2000, 4000}

var DefaultAbsorptionCoefficients = map[float64]float64{
	63:   0.10,
	125:  0.15,
	250:  0.20,
	500:  0.25,
	1000: 0.30,
	2000: 0.35,
	4000: 0.40,
}

func (r *RTEstimator) CalculateRT(surfaces []domain.Surface, materials map[string]*domain.Material) (domain.RTResults, []string) {
	var warnings []string
	var bands []domain.RTBand

	totalArea := 0.0
	for _, surface := range surfaces {
		totalArea += surface.AreaM2
	}

	for _, freq := range OctaveBands {
		totalAbsorption := 0.0
		hasMissingMaterial := false

		for _, surface := range surfaces {
			coeff := r.getAbsorptionCoefficient(&surface, freq, materials)
			if coeff < 0 {
				hasMissingMaterial = true
				coeff = DefaultAbsorptionCoefficients[freq]
			}
			totalAbsorption += coeff * surface.AreaM2
		}

		if hasMissingMaterial {
			warnings = append(warnings, "Missing material data at "+freqToString(freq)+" Hz, using fallback coefficients")
		}

		meanAbsorption := totalAbsorption / totalArea
		equivalentAbsorptionArea := totalAbsorption

		sabineT60 := r.calculateSabineT60(equivalentAbsorptionArea)
		eyringT60 := r.calculateEyringT60(equivalentAbsorptionArea, meanAbsorption)

		bands = append(bands, domain.RTBand{
			FrequencyHz:              freq,
			SabineT60:                sabineT60,
			EyringT60:                eyringT60,
			EquivalentAbsorptionArea: equivalentAbsorptionArea,
			MeanAbsorption:           meanAbsorption,
		})
	}

	return domain.RTResults{Bands: bands}, warnings
}

func (r *RTEstimator) getAbsorptionCoefficient(surface *domain.Surface, freq float64, materials map[string]*domain.Material) float64 {
	if surface.MaterialID == nil {
		return -1
	}

	material, ok := materials[*surface.MaterialID]
	if !ok {
		return -1
	}

	bands, err := material.GetAbsorptionBands()
	if err != nil {
		return -1
	}

	for _, band := range bands {
		if float64(band.Frequency) == freq {
			return band.Value
		}
	}

	return -1
}

func (r *RTEstimator) calculateSabineT60(equivalentAbsorptionArea float64) float64 {
	if equivalentAbsorptionArea <= 0.001 {
		return 0
	}
	return (0.161 * r.volume) / equivalentAbsorptionArea
}

func (r *RTEstimator) calculateEyringT60(equivalentAbsorptionArea, meanAbsorption float64) float64 {
	if meanAbsorption <= 0.001 || meanAbsorption >= 1.0 {
		return r.calculateSabineT60(equivalentAbsorptionArea)
	}
	surfaceArea := equivalentAbsorptionArea / meanAbsorption
	return (0.161 * r.volume) / (-surfaceArea * math.Log(1-meanAbsorption))
}

func freqToString(freq float64) string {
	if freq >= 1000 {
		return fmt.Sprintf("%.0fk", freq/1000)
	}
	return fmt.Sprintf("%.0f", freq)
}
