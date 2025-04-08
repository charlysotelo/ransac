package ransac

import (
	"iter"
	"math/rand"
	"time"

	"github.com/charlysotelo/ransac/internal"
)

// WithMaxIterations limits RANSAC to run only up to n iterations
func WithMaxIterations(n uint) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.terminationCondition |= MaxIterations
		p.maxIterations = n
	}
}

// WithMaxIterations limits RANSAC to run only up to duration d
func WithTimeout(d time.Duration) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.terminationCondition |= TimeLimit
		p.timeLimit = d
	}
}

// WithNoMaxIterations places no iteration limit on RANSAC
// In practice this means you must specify another termination condition
// such as WithTimeout
func WithNoMaxIterations() func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.terminationCondition = p.terminationCondition &^ MaxIterations
	}
}

// WithChooser sets the chooser function to be used for selecting random subsets of points
// The chooser function should yield combinations of indices for the points
// by default this uses the randomGonumChooser function
func WithChooser(chooser iter.Seq[[]int]) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.chooser = chooser
	}
}

// WithNumberOfWorkers sets the number of goroutines to be used as workers
// By default this is set to runtime.GOMAXPROCS(0) -- which should be your number of CPUs
func WithNumberOfWorkers(n uint) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.numberOfWorkers = n
	}
}

// WithNoConsensusSetFit enable or disables the consensus set fit typically done
// after the best model is found. This is enabled by default
func WithConsensusSetFit(doConsensusSetFit bool) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.doConsensusSetFit = doConsensusSetFit
	}
}

// WithRand sets the random number generator to be used for selecting random subsets of points
// This is useful for testing purposes
func WithRand(r *rand.Rand) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.localRand = r
	}
}

func applyOptions[R internal.Number, M Model[R]](p *problem[R, M], options ...any) {
	for _, opt := range options {
		switch opt := opt.(type) {
		case func(*problem[R, M]):
			opt(p)
		case func(*pNonGenericFields):
			opt(&(p.pNonGenericFields))
		default:
			panic("Invalid option type")
		}
	}
}
