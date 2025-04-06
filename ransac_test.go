package ransac

import (
	"testing"
)

func TestDefaultChooser(t *testing.T) {
	// Test case 1: Basic functionality
	N := uint(5)
	const K uint = uint(3)
	expectedCombinations := map[[K]uint]struct{}{
		[K]uint{0, 1, 2}: {},
		[K]uint{0, 1, 3}: {},
		[K]uint{0, 1, 4}: {},
		[K]uint{0, 2, 3}: {},
		[K]uint{0, 2, 4}: {},
		[K]uint{0, 3, 4}: {},
		[K]uint{1, 2, 3}: {},
		[K]uint{1, 2, 4}: {},
		[K]uint{1, 3, 4}: {},
		[K]uint{2, 3, 4}: {},
	}

	expectedCombinationsCount := len(expectedCombinations)
	indecesCount := 0
	for indeces := range defaultChooser(N, K) {
		indecesCount++
		asArray := [K]uint{}
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
