package ransac

import (
	"iter"
)

func defaultChooser(N, K uint) iter.Seq[[]uint] {
	// Generate all combinations of N choose K
	return func(yield func([]uint) bool) {
		comb := make([]uint, K)
		var generate func(int, int) bool
		generate = func(start, depth int) bool {
			if depth == int(K) {
				// Yield a copy of the current combination
				combCopy := make([]uint, K)
				copy(combCopy, comb)
				return yield(combCopy)
			}
			for i := start; i <= int(N)-int(K)+depth; i++ {
				comb[depth] = uint(i)
				if !generate(i+1, depth+1) {
					return false
				}
			}
			return true
		}
		generate(0, 0)
	}
}
