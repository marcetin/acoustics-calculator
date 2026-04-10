package acoustics

import (
	"math"

	"acoustics-calculator/internal/domain"
)

type ImageSourceSolver struct {
	width  float64
	length float64
	height float64
	c      float64
}

func NewImageSourceSolver(width, length, height, c float64) *ImageSourceSolver {
	return &ImageSourceSolver{
		width:  width,
		length: length,
		height: height,
		c:      c,
	}
}

func (is *ImageSourceSolver) CalculateReflections(sources []domain.Source, receivers []domain.Receiver, surfaces []domain.Surface) ([]domain.ReflectionPath, []string) {
	var warnings []string
	var reflections []domain.ReflectionPath

	surfaceMap := is.buildSurfaceMap(surfaces)

	for _, source := range sources {
		for _, receiver := range receivers {
			firstOrder := is.calculateFirstOrder(source, receiver, surfaceMap)
			reflections = append(reflections, firstOrder...)

			secondOrder := is.calculateSecondOrder(source, receiver, surfaceMap)
			if len(secondOrder) == 0 {
				warnings = append(warnings, "No valid second-order reflections for source "+source.Name+" to receiver "+receiver.Name)
			}
			reflections = append(reflections, secondOrder...)
		}
	}

	return reflections, warnings
}

func (is *ImageSourceSolver) buildSurfaceMap(surfaces []domain.Surface) map[string]domain.Surface {
	m := make(map[string]domain.Surface)
	for _, surface := range surfaces {
		m[surface.ID] = surface
	}
	return m
}

func (is *ImageSourceSolver) calculateFirstOrder(source domain.Source, receiver domain.Receiver, surfaceMap map[string]domain.Surface) []domain.ReflectionPath {
	var paths []domain.ReflectionPath

	imagePositions := []struct {
		x, y, z float64
		surface string
	}{
		{-source.PositionX, source.PositionY, source.PositionZ, "left"},
		{2*is.width - source.PositionX, source.PositionY, source.PositionZ, "right"},
		{source.PositionX, -source.PositionY, source.PositionZ, "front"},
		{source.PositionX, 2*is.length - source.PositionY, source.PositionZ, "back"},
		{source.PositionX, source.PositionY, -source.PositionZ, "floor"},
		{source.PositionX, source.PositionY, 2*is.height - source.PositionZ, "ceiling"},
	}

	for _, img := range imagePositions {
		pathLength := math.Sqrt(
			math.Pow(img.x-receiver.PositionX, 2) +
				math.Pow(img.y-receiver.PositionY, 2) +
				math.Pow(img.z-receiver.PositionZ, 2),
		)

		arrivalTimeMs := (pathLength / is.c) * 1000

		directPathLength := math.Sqrt(
			math.Pow(source.PositionX-receiver.PositionX, 2) +
				math.Pow(source.PositionY-receiver.PositionY, 2) +
				math.Pow(source.PositionZ-receiver.PositionZ, 2),
		)

		if pathLength <= directPathLength {
			continue
		}

		energy := is.calculateEnergy(pathLength)

		surfaceID := is.findSurfaceIDByKind(img.surface, surfaceMap)

		paths = append(paths, domain.ReflectionPath{
			SourceID:        source.ID,
			ReceiverID:      receiver.ID,
			Order:           1,
			SurfaceSequence: []string{surfaceID},
			PathLengthM:     pathLength,
			ArrivalTimeMs:   arrivalTimeMs,
			EnergyEstimate:  energy,
		})
	}

	return paths
}

func (is *ImageSourceSolver) calculateSecondOrder(source domain.Source, receiver domain.Receiver, surfaceMap map[string]domain.Surface) []domain.ReflectionPath {
	var paths []domain.ReflectionPath

	surfaces := []string{"left", "right", "front", "back", "floor", "ceiling"}

	for i, s1 := range surfaces {
		for _, s2 := range surfaces[i+1:] {
			img1X, img1Y, img1Z := is.getImageSource(source, s1)
			img2X, img2Y, img2Z := is.getImageSourceFromImage(img1X, img1Y, img1Z, s2)

			pathLength := math.Sqrt(
				math.Pow(img2X-receiver.PositionX, 2) +
					math.Pow(img2Y-receiver.PositionY, 2) +
					math.Pow(img2Z-receiver.PositionZ, 2),
			)

			directPathLength := math.Sqrt(
				math.Pow(source.PositionX-receiver.PositionX, 2) +
					math.Pow(source.PositionY-receiver.PositionY, 2) +
					math.Pow(source.PositionZ-receiver.PositionZ, 2),
			)

			if pathLength <= directPathLength {
				continue
			}

			arrivalTimeMs := (pathLength / is.c) * 1000
			energy := is.calculateEnergy(pathLength)

			surfaceID1 := is.findSurfaceIDByKind(s1, surfaceMap)
			surfaceID2 := is.findSurfaceIDByKind(s2, surfaceMap)

			if surfaceID1 == "" || surfaceID2 == "" {
				continue
			}

			paths = append(paths, domain.ReflectionPath{
				SourceID:        source.ID,
				ReceiverID:      receiver.ID,
				Order:           2,
				SurfaceSequence: []string{surfaceID1, surfaceID2},
				PathLengthM:     pathLength,
				ArrivalTimeMs:   arrivalTimeMs,
				EnergyEstimate:  energy,
			})
		}
	}

	return paths
}

func (is *ImageSourceSolver) getImageSource(source domain.Source, surface string) (x, y, z float64) {
	switch surface {
	case "left":
		return -source.PositionX, source.PositionY, source.PositionZ
	case "right":
		return 2*is.width - source.PositionX, source.PositionY, source.PositionZ
	case "front":
		return source.PositionX, -source.PositionY, source.PositionZ
	case "back":
		return source.PositionX, 2*is.length - source.PositionY, source.PositionZ
	case "floor":
		return source.PositionX, source.PositionY, -source.PositionZ
	case "ceiling":
		return source.PositionX, source.PositionY, 2*is.height - source.PositionZ
	default:
		return source.PositionX, source.PositionY, source.PositionZ
	}
}

func (is *ImageSourceSolver) getImageSourceFromImage(x, y, z float64, surface string) (float64, float64, float64) {
	switch surface {
	case "left":
		return -x, y, z
	case "right":
		return 2*is.width - x, y, z
	case "front":
		return x, -y, z
	case "back":
		return x, 2*is.length - y, z
	case "floor":
		return x, y, -z
	case "ceiling":
		return x, y, 2*is.height - z
	default:
		return x, y, z
	}
}

func (is *ImageSourceSolver) calculateEnergy(pathLength float64) float64 {
	if pathLength <= 0 {
		return 0
	}
	return 1.0 / (pathLength * pathLength)
}

func (is *ImageSourceSolver) findSurfaceIDByKind(kind string, surfaceMap map[string]domain.Surface) string {
	for _, surface := range surfaceMap {
		surfaceKind := string(surface.Kind)
		if surfaceKind == kind {
			return surface.ID
		}
	}
	return ""
}
