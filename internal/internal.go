package internal

import "iter"

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// OrderedChooser generates all combinations of N choose K
// It is not random, but it is deterministic. Useful for testing
func OrderedChooser(N, K uint) iter.Seq[[]int] {
	// Generate all combinations of N choose K
	return func(yield func([]int) bool) {
		comb := make([]int, K)
		var generate func(int, int) bool
		generate = func(start, depth int) bool {
			if depth == int(K) {
				// Yield a copy of the current combination
				combCopy := make([]int, K)
				copy(combCopy, comb)
				return yield(combCopy)
			}
			for i := start; i <= int(N)-int(K)+depth; i++ {
				comb[depth] = int(i)
				if !generate(i+1, depth+1) {
					return false
				}
			}
			return true
		}
		generate(0, 0)
	}
}
