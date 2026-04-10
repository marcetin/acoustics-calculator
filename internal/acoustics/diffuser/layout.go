package diffuser

import (
	"math"

	"acoustics-calculator/internal/domain"
)

type LayoutSynthesizer struct{}

func NewLayoutSynthesizer() *LayoutSynthesizer {
	return &LayoutSynthesizer{}
}

func (l *LayoutSynthesizer) SynthesizeLayout(patch *domain.BoundingBox, diffuser *domain.DiffuserType, orientation string) *domain.LayoutResult {
	panelWidth := diffuser.PanelWidthM
	panelHeight := diffuser.PanelHeightM

	patchWidth := patch.MaxX - patch.MinX
	patchHeight := patch.MaxY - patch.MinY

	if patchWidth <= 0 || patchHeight <= 0 {
		return &domain.LayoutResult{
			BoundingBox:    patch,
			PanelCenters:   []domain.Vec3{},
			PanelCount:     0,
			Orientation:    orientation,
			CoverageAreaM2: 0,
		}
	}

	cols := int(math.Floor(patchWidth / panelWidth))
	rows := int(math.Floor(patchHeight / panelHeight))

	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	var centers []domain.Vec3
	startX := patch.MinX + (patchWidth - float64(cols)*panelWidth)/2 + panelWidth/2
	startY := patch.MinY + (patchHeight - float64(rows)*panelHeight)/2 + panelHeight/2
	z := (patch.MinZ + patch.MaxZ) / 2

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x := startX + float64(col)*panelWidth
			y := startY + float64(row)*panelHeight
			centers = append(centers, domain.Vec3{X: x, Y: y, Z: z})
		}
	}

	panelCount := cols * rows
	coverageAreaM2 := float64(panelCount) * panelWidth * panelHeight

	return &domain.LayoutResult{
		BoundingBox:    patch,
		PanelCenters:   centers,
		PanelCount:     panelCount,
		Orientation:    orientation,
		CoverageAreaM2: coverageAreaM2,
	}
}
