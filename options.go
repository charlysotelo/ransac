package ransac

import (
	"iter"
	"time"

	"github.com/charlysotelo/ransac/internal"
)

func WithMaxIterations(n uint) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.terminationCondition |= MaxIterations
		p.maxIterations = n
	}
}

func WithTimeout(d time.Duration) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.terminationCondition |= TimeLimit
		p.timeLimit = d
	}
}

func WithExhaustedIterations() func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.terminationCondition = p.terminationCondition &^ MaxIterations
	}
}

// WithChooser sets the chooser function to be used for selecting random subsets of points.
// The chooser function should yield combinations of indices for the points.
func WithChooser(chooser iter.Seq[[]uint]) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.chooser = chooser
	}
}

func WithNumberOfWorkers(n uint) func(*pNonGenericFields) {
	return func(p *pNonGenericFields) {
		p.numberOfWorkers = n
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
