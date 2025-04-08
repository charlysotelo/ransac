package ransac

import (
	"iter"
	"math/rand"

	"gonum.org/v1/gonum/stat/combin"
)

func defaultChooser(N, K uint) iter.Seq[[]int] {
	return RandomGonumChooser(N, K)
}

// RandomGonumChooser generates random combinations of N choose K
// using gonum's combin.IndexToCombination function with a random index
// Note that it may yield the same combination multiple times
func RandomGonumChooser(N, K uint) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		for {

			result := make([]int, K)
			combin.IndexToCombination(result, rand.Intn(int(N)), int(N), int(K))

			if !yield(result) {
				return
			}
		}
	}
}
