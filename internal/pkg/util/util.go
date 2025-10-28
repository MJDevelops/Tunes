package util

import (
	"math"
	"math/rand/v2"
)

func FisherYates[T any](s []T) []T {
	for i := len(s) - 1; i >= 1; i-- {
		j := int(math.Floor(rand.Float64() * (float64(i) + 1)))
		s[i], s[j] = s[j], s[i]
	}

	return s
}
