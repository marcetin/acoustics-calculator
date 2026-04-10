package diffuser

import (
	"acoustics-calculator/internal/domain"
)

type Catalog struct {
	diffusers []domain.DiffuserType
}

func NewCatalog(diffusers []*domain.DiffuserType) *Catalog {
	c := &Catalog{}
	for _, d := range diffusers {
		c.diffusers = append(c.diffusers, *d)
	}
	return c
}

func (c *Catalog) FindCompatible(targetBands []float64, maxDepthM float64, patchAreaM2 float64, orientation string) *domain.DiffuserType {
	var bestMatch *domain.DiffuserType
	bestScore := 0.0

	for _, d := range c.diffusers {
		if d.PanelDepthM > maxDepthM {
			continue
		}

		panelArea := d.PanelWidthM * d.PanelHeightM
		if panelArea > patchAreaM2 {
			continue
		}

		score := c.calculateCompatibilityScore(&d, targetBands, orientation)
		if score > bestScore {
			bestScore = score
			bestMatch = &d
		}
	}

	return bestMatch
}

func (c *Catalog) calculateCompatibilityScore(diffuser *domain.DiffuserType, targetBands []float64, orientation string) float64 {
	score := 0.0

	bandMatch := 0.0
	for _, band := range targetBands {
		if band >= diffuser.MinEffectiveHz && band <= diffuser.MaxEffectiveHz {
			bandMatch += 1.0
		}
	}
	if len(targetBands) > 0 {
		bandMatch = bandMatch / float64(len(targetBands))
	}
	score += bandMatch * 0.5

	if diffuser.Category == domain.DiffuserCategoryQRD2D {
		score += 0.3
	} else if diffuser.Category == domain.DiffuserCategoryQRD1D {
		score += 0.2
	}

	panelArea := diffuser.PanelWidthM * diffuser.PanelHeightM
	if panelArea >= 0.36 {
		score += 0.2
	}

	return score
}

func (c *Catalog) GetByCategory(category domain.DiffuserCategory) []domain.DiffuserType {
	var result []domain.DiffuserType
	for _, d := range c.diffusers {
		if d.Category == category {
			result = append(result, d)
		}
	}
	return result
}
