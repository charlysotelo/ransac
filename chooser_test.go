package ransac

import (
	"testing"

	"github.com/charlysotelo/ransac/internal"
)

func TestOrderedChooser(t *testing.T) {
	// Test case 1: Basic functionality
	N := uint(5)
	const K uint = uint(3)
	expectedCombinations := map[[K]int]struct{}{
		{0, 1, 2}: {},
		{0, 1, 3}: {},
		{0, 1, 4}: {},
		{0, 2, 3}: {},
		{0, 2, 4}: {},
		{0, 3, 4}: {},
		{1, 2, 3}: {},
		{1, 2, 4}: {},
		{1, 3, 4}: {},
		{2, 3, 4}: {},
	}

	expectedCombinationsCount := len(expectedCombinations)
	indecesCount := 0
	for indeces := range internal.OrderedChooser(N, K) {
		indecesCount++
		asArray := [K]int{}
		copy(asArray[:], indeces)
		if _, ok := expectedCombinations[asArray]; !ok {
			t.Errorf("Unexpected combination: %v", indeces)
		}
		delete(expectedCombinations, asArray)
	}

	if indecesCount != expectedCombinationsCount {
		t.Errorf("Expected %d combinations, got %d", expectedCombinationsCount, indecesCount)
	}
}
