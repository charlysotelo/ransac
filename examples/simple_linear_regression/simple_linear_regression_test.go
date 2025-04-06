package simple_linear_regression

import (
	"fmt"

	"github.com/charlysotelo/ransac"
)

func ExampleLinearRegressionModel_ransac() {
	// Example data points
	points := [][]float64{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{16, 0},
		{-37, 43},
	}

	// Create a new LinearRegressionModel with a threshold of 0.5
	model := NewLinearRegressionModel(0.5)

	// Fit the model to the data points
	ransac.ModelFit(points, model,
		// Optional: Customize the RANSAC algorithm parameters
		ransac.WithMaxIterations(1000),
		// ransac.WithExhaustedIterations(),
		// ransac.WithTimeout(5 * time.Second),
		// ransac.WithNumberOfWorkers(4),
	)

	// At this point your model is fitted to the data points
	fmt.Println(model)

	// And may be used to classify points as inliers or outliers
	for _, point := range points {
		if model.IsInlier(point) {
			fmt.Printf("Point %v is an inlier\n", point)
			continue
		}
		fmt.Printf("Point %v is an outlier\n", point)
	}

	// Output:
	// y = 1.00x + 1.00
	// Point [0 1] is an inlier
	// Point [1 2] is an inlier
	// Point [2 3] is an inlier
	// Point [3 4] is an inlier
	// Point [16 0] is an outlier
	// Point [-37 43] is an outlier
}
