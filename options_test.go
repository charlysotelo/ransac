package ransac

import (
	"math/rand"
	"testing"
	"time"
)

type dummyModel struct{}

func (dummyModel) Fit([][]float64) (dummyModel, error) { return dummyModel{}, nil }
func (dummyModel) Error([]float64) float64             { return 0 }
func (dummyModel) Copy() dummyModel                    { return dummyModel{} }

func TestWithMaxIterations(t *testing.T) {
	p := &pNonGenericFields{}
	f := WithMaxIterations(42)
	f(p)
	if p.maxIterations != 42 {
		t.Errorf("expected maxIterations=42, got %d", p.maxIterations)
	}
	if p.terminationCondition&MaxIterations == 0 {
		t.Error("expected MaxIterations bit set in terminationCondition")
	}
}

func TestWithTimeout(t *testing.T) {
	p := &pNonGenericFields{}
	d := 2 * time.Second
	f := WithTimeout(d)
	f(p)
	if p.timeLimit != d {
		t.Errorf("expected timeLimit=%v, got %v", d, p.timeLimit)
	}
	if p.terminationCondition&TimeLimit == 0 {
		t.Error("expected TimeLimit bit set in terminationCondition")
	}
}

func TestWithNoMaxIterations(t *testing.T) {
	p := &pNonGenericFields{terminationCondition: MaxIterations | TimeLimit}
	f := WithNoMaxIterations()
	f(p)
	if p.terminationCondition&MaxIterations != 0 {
		t.Error("expected MaxIterations bit cleared in terminationCondition")
	}
}

func TestWithChooser(t *testing.T) {
	p := &pNonGenericFields{}
	called := false
	chooser := func(yield func([]int) bool) {
		called = true
	}
	f := WithChooser(chooser)
	f(p)
	if p.chooser == nil {
		t.Error("expected chooser to be set")
	}
	p.chooser(func([]int) bool { return false })
	if !called {
		t.Error("expected chooser to be called")
	}
}

func TestWithNumberOfWorkers(t *testing.T) {
	p := &pNonGenericFields{}
	f := WithNumberOfWorkers(7)
	f(p)
	if p.numberOfWorkers != 7 {
		t.Errorf("expected numberOfWorkers=7, got %d", p.numberOfWorkers)
	}
}

func TestWithConsensusSetFit(t *testing.T) {
	p := &pNonGenericFields{}
	f := WithConsensusSetFit(false)
	f(p)
	if p.doConsensusSetFit != false {
		t.Errorf("expected doConsensusSetFit=false, got %v", p.doConsensusSetFit)
	}
	f2 := WithConsensusSetFit(true)
	f2(p)
	if p.doConsensusSetFit != true {
		t.Errorf("expected doConsensusSetFit=true, got %v", p.doConsensusSetFit)
	}
}

func TestWithRand(t *testing.T) {
	p := &pNonGenericFields{}
	r := rand.New(rand.NewSource(1))
	f := WithRand(r)
	f(p)
	if p.localRand != r {
		t.Error("expected localRand to be set")
	}
}
