package simple_linear_regression

import (
	"fmt"
	"math"
)

// LinearRegressionModel implements the Model interface for a simple 2D linear regression.
type LinearRegressionModel struct {
	slope     float64
	intercept float64
	threshold float64 // Threshold to determine if a point is an inlier
}

// MinimalFitpoints returns the minimum number of points needed to fit a line (2 points).
func (m *LinearRegressionModel) MinimalFitpoints() uint {
	return 2
}

// Fit calculates the slope and intercept of the line using the given points
// by performing a linear least-squares regression.
func (m *LinearRegressionModel) Fit(points [][]float64) {
	if len(points) < 2 {
		panic("LinearRegressionModel requires at least 2 points to fit")
	}

	var sumX, sumY, sumXY, sumX2 float64
	n := float64(len(points))

	for _, point := range points {
		if len(point) != 2 {
			panic("Each point must have exactly 2 dimensions")
		}
		x, y := point[0], point[1]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope and intercept using least-squares formulas
	m.slope = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	m.intercept = (sumY - m.slope*sumX) / n
}

// IsInlier checks if a point lies close enough to the fitted line (within the threshold).
func (m *LinearRegressionModel) IsInlier(point []float64) bool {
	if len(point) != 2 {
		panic("Point must have exactly 2 dimensions")
	}

	x, y := point[0], point[1]
	// Calculate the distance from the point to the line
	yPredicted := m.slope*x + m.intercept
	distance := math.Abs(y - yPredicted)

	return distance <= m.threshold
}

func (m *LinearRegressionModel) Copy() *LinearRegressionModel {
	result := *m
	return &result
}

func (m *LinearRegressionModel) String() string {
	return fmt.Sprintf("y = %.2fx + %.2f", m.slope, m.intercept)
}

// NewLinearRegressionModel creates a new LinearRegressionModel with the given inlier threshold.
func NewLinearRegressionModel(threshold float64) *LinearRegressionModel {
	return &LinearRegressionModel{
		threshold: threshold,
	}
}
