package ransac

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

type testModel struct {
	fitCalled      int
	isInlierCalled int
	minFit         uint
	fitPoints      [][]float64
	inlierMask     []bool
}

func (m *testModel) MinimalFitpoints() uint { return m.minFit }
func (m *testModel) Fit(points [][]float64) {
	m.fitCalled++
	m.fitPoints = points
}
func (m *testModel) IsInlier(point []float64) bool {
	m.isInlierCalled++
	if m.inlierMask == nil {
		return true
	}
	// Assume point[0] is index for test
	idx := int(point[0])
	if idx < len(m.inlierMask) {
		return m.inlierMask[idx]
	}
	return false
}
func (m *testModel) Copy() *testModel {
	copy := *m
	return &copy
}

func makeData(n, d int) [][]float64 {
	data := make([][]float64, n)
	for i := range data {
		row := make([]float64, d)
		row[0] = float64(i) // for inlier mask
		data[i] = row
	}
	return data
}

func TestModelFit_Basic(t *testing.T) {
	data := makeData(10, 2)
	model := &testModel{minFit: 2}
	err := ModelFit(data, model)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if model.fitCalled == 0 {
		t.Error("expected Fit to be called at least once")
	}
}

func TestModelFit_NoModelCopierError(t *testing.T) {
	data := makeData(10, 2)
	type noCopyModel struct{ testModel }
	model := &noCopyModel{testModel{minFit: 2}}
	// Remove CopyableModel interface by not embedding Copy()
	// Use >1 worker to trigger error
	err := ModelFit(data, model, WithNumberOfWorkers(2))
	if !errors.Is(err, ErrNoModelCopier) {
		t.Errorf("expected ErrNoModelCopier, got %v", err)
	}
}

func TestModelFit_InvalidParams(t *testing.T) {
	data := makeData(10, 2)
	model := &testModel{minFit: 2}
	err := ModelFit(data, model, WithNumberOfWorkers(0))
	if err == nil || err.Error() == "" {
		t.Error("expected error for invalid parameters")
	}
}

func TestModelFit_ConsensusSetFitFalse(t *testing.T) {
	data := makeData(10, 2)
	model := &testModel{minFit: 2}
	err := ModelFit(data, model, WithConsensusSetFit(false))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestModelFit_WithTimeout(t *testing.T) {
	data := makeData(10, 2)
	model := &testModel{minFit: 2}
	err := ModelFit(data, model, WithTimeout(1*time.Millisecond))
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected no error or context.DeadlineExceeded, got %v", err)
	}
}

func TestProblem_validateParams(t *testing.T) {
	data := makeData(5, 2)
	model := &testModel{minFit: 2}
	p := getDefaultProblem(data, model)
	p.minInliers = 0
	if err := p.validateParams(); !errors.Is(err, ErrMinInliersZero) {
		t.Errorf("expected ErrMinInliersZero, got %v", err)
	}
	p.minInliers = 10
	if err := p.validateParams(); !errors.Is(err, ErrMinInliers) {
		t.Errorf("expected ErrMinInliers, got %v", err)
	}
	p.minInliers = 2
	p.numberOfWorkers = 0
	if err := p.validateParams(); !errors.Is(err, ErrNumWorkersZero) {
		t.Errorf("expected ErrNumWorkersZero, got %v", err)
	}
	p.numberOfWorkers = 1
	p.terminationCondition = MaxIterations
	p.maxIterations = 0
	if err := p.validateParams(); !errors.Is(err, ErrMaxIterationsZero) {
		t.Errorf("expected ErrMaxIterationsZero, got %v", err)
	}
}

func TestSelectSubset(t *testing.T) {
	set := [][]float64{{1, 2}, {3, 4}, {5, 6}}
	subset := make([][]float64, 2)
	for i := range subset {
		subset[i] = make([]float64, 2)
	}
	selectSubset(set, subset, []int{2, 0})
	if subset[0][0] != 5 || subset[1][0] != 1 {
		t.Errorf("unexpected subset: got %v", subset)
	}
}

func TestGetDefaultProblem(t *testing.T) {
	data := makeData(3, 2)
	model := &testModel{minFit: 2}
	p := getDefaultProblem(data, model)
	if p.maxIterations != 1000 {
		t.Errorf("expected maxIterations=1000, got %d", p.maxIterations)
	}
	if p.minInliers != 2 {
		t.Errorf("expected minInliers=2, got %d", p.minInliers)
	}
	if p.numberOfWorkers != uint(runtime.GOMAXPROCS(0)) {
		t.Errorf("expected numberOfWorkers=%d, got %d", runtime.GOMAXPROCS(0), p.numberOfWorkers)
	}
	if p.chooser == nil {
		t.Error("expected chooser to be set")
	}
	if p.data == nil || p.model == nil {
		t.Error("expected data and model to be set")
	}
	if !p.doConsensusSetFit {
		t.Error("expected doConsensusSetFit to be true")
	}
	if p.localRand == nil {
		t.Error("expected localRand to be set")
	}
}
