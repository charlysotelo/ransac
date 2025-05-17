package ransac

import (
	"testing"
	"time"
)

type behaviorModel struct {
	fitCount      int
	consensusMask []bool
	minFit        uint
}

func (m *behaviorModel) MinimalFitpoints() uint { return m.minFit }
func (m *behaviorModel) Fit(points [][]float64) { m.fitCount++ }
func (m *behaviorModel) IsInlier(point []float64) bool {
	if m.consensusMask == nil {
		return true
	}
	idx := int(point[0])
	if idx < len(m.consensusMask) {
		return m.consensusMask[idx]
	}
	return false
}
func (m *behaviorModel) Copy() *behaviorModel {
	copy := *m
	return &copy
}

func makeBehaviorData(n int) [][]float64 {
	data := make([][]float64, n)
	for i := range data {
		row := make([]float64, 1)
		row[0] = float64(i)
		data[i] = row
	}
	return data
}

func TestWithMaxIterations_behavior(t *testing.T) {
	data := makeBehaviorData(10)
	model := &behaviorModel{minFit: 2}
	ModelFit(data, model, WithMaxIterations(3))
	if model.fitCount > 3 {
		t.Errorf("expected fitCount <= 3, got %d", model.fitCount)
	}
}

func TestWithTimeout_behavior(t *testing.T) {
	data := makeBehaviorData(10)
	model := &behaviorModel{minFit: 2}
	start := time.Now()
	ModelFit(data, model, WithTimeout(10*time.Millisecond))
	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("expected ModelFit to return quickly, took %v", elapsed)
	}
}

func TestWithChooser_behavior(t *testing.T) {
	data := makeBehaviorData(5)
	model := &behaviorModel{minFit: 2}
	calls := 0
	chooser := func(yield func([]int) bool) {
		calls++
		yield([]int{0, 1})
		yield([]int{2, 3})
	}
	ModelFit(data, model, WithChooser(chooser))
	if calls == 0 {
		t.Error("expected chooser to be called")
	}
	if model.fitCount < 2 {
		t.Errorf("expected fitCount >= 2, got %d", model.fitCount)
	}
}

func TestWithConsensusSetFit_behavior(t *testing.T) {
	data := makeBehaviorData(5)
	model := &behaviorModel{minFit: 2}
	ModelFit(data, model, WithConsensusSetFit(false))
	// Should not error, and fitCount should be at least 1
	if model.fitCount < 1 {
		t.Errorf("expected fitCount >= 1, got %d", model.fitCount)
	}
}

func TestWithRand_behavior(t *testing.T) {
	// This test just ensures WithRand does not panic and can be set
	data := makeBehaviorData(5)
	model := &behaviorModel{minFit: 2}
	ModelFit(data, model, WithRand(nil))
	if model.fitCount < 1 {
		t.Errorf("expected fitCount >= 1, got %d", model.fitCount)
	}
}
