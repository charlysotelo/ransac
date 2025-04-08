// Package ransac implements random sample consensus (RANSAC)
//
// See https://en.wikipedia.org/wiki/Random_sample_consensus
package ransac

import (
	"context"
	"fmt"
	"iter"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/charlysotelo/ransac/internal"
)

// Model is RANSAC's interface to your model implementation
// The type parameter R is your number type for your model (typically float64)
type Model[R internal.Number] interface {
	MinimalFitpoints() uint // Minimal number of points needed to fit the model (e.g. 2 for line)
	Fit([][]R)              // Modifies internal state of model
	IsInlier([]R) bool      // Returns true if the point is an inlier
}

// See CopyableModel
type Copier[T any] interface {
	Copy() T
}

// CopyableModel represents a Model with a Copy() method. Since this
// RASNAC implementation modifies the internal state of a Model when it is fitted,
// each goroutine requires its own distinct model copy
type CopyableModel[R internal.Number, M Model[R]] interface {
	Model[R]
	Copier[M]
}

type terminationCondition int

const (
	// Termination condition
	MaxIterations terminationCondition = 1 << iota // Maximum number of iterations
	TimeLimit
)

type (
	problem[R internal.Number, M Model[R]] struct {
		pNonGenericFields
		pGenericFields[R, M]
	}

	pNonGenericFields struct {
		terminationCondition terminationCondition // Termination condition
		maxIterations        uint                 // Maximum number of iterations
		minInliers           uint                 // Minimum number of inliers to accept the model
		timeLimit            time.Duration        // Time limit for the algorithm
		chooser              iter.Seq[[]int]      // Chooses a random subset of points
		numberOfWorkers      uint                 // Number of workers to use for parallel processing
		doConsensusSetFit    bool                 // Whether to fit the model to the consensus set
		localRand            *rand.Rand           // useful for overrides in tests
	}

	pGenericFields[R internal.Number, M Model[R]] struct {
		data  [][]R // Data points
		model M
	}
)

func (p *problem[R, M]) shouldTerminate(iterationCount uint) bool {
	if p.terminationCondition&MaxIterations == 0 {
		return false
	}

	if iterationCount < p.maxIterations {
		return true
	}
	return false
}

func selectSubset[R internal.Number](set [][]R, subset [][]R, indeces []int) {
	i := 0
	for _, index := range indeces {
		subset[i] = set[index]
		i++
	}
}

type workerResult struct {
	indeces []int
	score   uint
}

func (p *problem[R, M]) worker(model Model[R], jobs <-chan []int, results chan<- workerResult, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	var result workerResult
	// Pre-allocate memory for hypothetical inliers
	hypotheticalInliers := make([][]R, model.MinimalFitpoints())
	for i := range hypotheticalInliers {
		hypotheticalInliers[i] = make([]R, len(p.data[0]))
	}

	for {

		// Timeout/Cancel check
		// Split into two so we always check for cancelation
		select {
		case <-ctx.Done():
			results <- result
			return
		default:
		}

		var indeces []int
		var ok bool
		select {
		case <-ctx.Done():
			results <- result
			return
		case indeces, ok = <-jobs:
			if !ok {
				results <- result
				return
			}
		}

		// for indeces := range jobs {
		// Fit model to hypothetical inliers
		selectSubset(p.data, hypotheticalInliers, indeces)
		model.Fit(hypotheticalInliers)

		// Count members in consensus set
		consensus := uint(0)
		for _, point := range p.data {
			if model.IsInlier(point) {
				consensus++
			}
		}

		// Keep track of the best hypothesis
		if consensus >= result.score {
			result.score = consensus
			result.indeces = indeces
		}
	}
}

func (p *problem[R, M]) modelFit(model CopyableModel[R, M]) error {
	// Create worker pool
	jobs := make(chan []int, p.numberOfWorkers)
	results := make(chan workerResult, p.numberOfWorkers)

	ctx, cancel := context.WithCancel(context.Background())
	if p.terminationCondition&TimeLimit != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), p.timeLimit)
	}
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(int(p.numberOfWorkers))
	for range p.numberOfWorkers {
		go p.worker(model.Copy(), jobs, results, ctx, &wg)
	}

	iterationCount := uint(0)
	for hypotheticalInliersIndeces := range p.chooser {
		jobs <- hypotheticalInliersIndeces
		iterationCount++

		if p.shouldTerminate(iterationCount) {
			break
		}
	}
	close(jobs)
	wg.Wait()
	close(results)

	// Gather best result
	var bestHypothesisIndeces []int
	bestHypothesisScore := uint(0)
	for result := range results {
		if result.score > bestHypothesisScore {
			bestHypothesisScore = result.score
			bestHypothesisIndeces = result.indeces
		}
	}

	// Fit the model to the best hypothesis
	bestHypothesis := make([][]R, model.MinimalFitpoints())
	for i := range len(bestHypothesis) {
		bestHypothesis[i] = make([]R, len(p.data[0]))
	}
	selectSubset(p.data, bestHypothesis, bestHypothesisIndeces)
	model.Fit(bestHypothesis)

	if !p.doConsensusSetFit {
		return nil
	}

	// Collect the consensus set
	consensusSet := make([][]R, 0, bestHypothesisScore)
	for _, point := range p.data {
		if model.IsInlier(point) {
			consensusSet = append(consensusSet, point)
		}
	}

	// Fit the model to consensus set
	model.Fit(consensusSet)
	return nil
}

// copyStub is a stub for the ModelCopier interface
// It is used when the model does not implement ModelCopier
// It simply returns the original model when Copy() is called
type copyStub[R internal.Number, M Model[R]] struct {
	Model M
}

func (c *copyStub[R, M]) MinimalFitpoints() uint {
	return c.Model.MinimalFitpoints()
}

func (c *copyStub[R, M]) Fit(points [][]R) {
	c.Model.Fit(points)
}
func (c *copyStub[R, M]) IsInlier(point []R) bool {
	return c.Model.IsInlier(point)
}
func (c *copyStub[R, M]) Copy() M {
	return c.Model
}

func (p *problem[R, M]) validateParams() error {
	if p.minInliers == 0 {
		return ErrMinInliersZero
	}
	if p.minInliers > uint(len(p.data)) {
		return ErrMinInliers
	}
	if p.numberOfWorkers == 0 {
		return ErrNumWorkersZero
	}
	if p.terminationCondition == MaxIterations && p.maxIterations == 0 {
		return ErrMaxIterationsZero
	}
	return nil
}

var pkgRand *rand.Rand

func init() {
	pkgRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getDefaultProblem[R internal.Number, M Model[R]](data [][]R, model M) *problem[R, M] {
	p := &problem[R, M]{}
	p.terminationCondition = MaxIterations
	p.maxIterations = 1000
	p.minInliers = model.MinimalFitpoints()
	p.numberOfWorkers = uint(runtime.GOMAXPROCS(0))
	p.chooser = defaultChooser(uint(len(data)), p.minInliers)
	p.data = data
	p.model = model
	p.doConsensusSetFit = true
	p.localRand = pkgRand
	return p
}

// ModelFit fits the model to the supplied data using RANSAC
func ModelFit[R internal.Number, M Model[R]](data [][]R, model M, options ...any) error {
	p := getDefaultProblem(data, model)
	applyOptions(p, options...)

	if err := p.validateParams(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Can't use type-assertions in structs so we wrap in any()
	if copier, ok := any(model).(CopyableModel[R, M]); ok {
		return p.modelFit(copier)
	}

	if p.numberOfWorkers > 1 {
		return ErrNoModelCopier
	}

	return p.modelFit(&copyStub[R, M]{model})
}
