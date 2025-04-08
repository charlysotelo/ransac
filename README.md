# RANSAC Go

[![Build Status](https://github.com/charlysotelo/ransac/actions/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/charlysotelo/ransac/actions/workflows/tests.yaml?query=branch%3Amain)
[![go.dev reference](https://pkg.go.dev/badge/github.com/charlysotelo/ransac)](https://pkg.go.dev/github.com/charlysotelo/ransac)
[![Go Report Card](https://goreportcard.com/badge/github.com/charlysotelo/ransac)](https://goreportcard.com/report/github.com/charlysotelo/ransac)
[![codecov](https://codecov.io/github/charlysotelo/ransac/graph/badge.svg?token=118ZWPMA3G)](https://codecov.io/github/charlysotelo/ransac)


This is a golang implementation of [RANSAC](https://en.wikipedia.org/wiki/Random_sample_consensus). At a minimum, you provide a model which implements the `Model` interface. See [simple_linear_regression.go](examples/simple_linear_regression/simple_linear_regression.go) for an example implementation and [simple_linear_regression_test.go](examples/simple_linear_regression/simple_linear_regression_test.go) for its usage

## Usage
```go
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

	// Create a new LinearRegressionModel, which is an implementation of ransac.Model,
	// with a threshold of 0.5
	// note you provide an implementation of Model -- NewLinearRegressionModel
	// is provided by this module for demonstration purposes
	model := NewLinearRegressionModel(0.5)

	// Fit the model to the data points
	ransac.ModelFit(points, model,
		// Optional: Customize the RANSAC algorithm parameters
		// ransac.WithMaxIterations(1000),
		// ransac.WithTimeout(5 * time.Second),
		// ransac.WithNumberOfWorkers(4),
		// ransac.WithConsensusSetFit(true),
		// ransac.WithRand(...),
		// ... see options.go for more options
		//
		// This chooser was chosen for deterministic behavior during testing
		ransac.WithChooser(internal.OrderedChooser(uint(len(points)), model.MinimalFitpoints())),
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
```

## A note on efficiency
RANSAC is highly parallelizable, but to make full use of it here your model needs to also implement the `Copier` interface. Again, see [simple_linear_regression.go](examples/simple_linear_regression/simple_linear_regression.go) for an example
