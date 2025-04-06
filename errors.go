package ransac

import "errors"

var (
	ErrNoModelCopier     = errors.New("Model does not implement ModelCopier, cannot use multiple workers")
	ErrMinInliersZero    = errors.New("minimum inliers cannot be zero")
	ErrMinInliers        = errors.New("minimum inliers cannot be greater than the number of data points")
	ErrNumWorkersZero    = errors.New("number of workers cannot be zero")
	ErrMaxIterationsZero = errors.New("max iterations cannot be zero")
)
