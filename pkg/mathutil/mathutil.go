package mathutil

import (
	"math"
	"sort"
)

func Percentile(percentile float64, input []int) int {
	if len(input) == 0 {
		return 0
	}

	if len(input) == 1 {
		return input[0]
	}

	sort.Ints(input)
	index := int(math.Round(percentile * float64(len(input))))
	if index == len(input) {
		return input[index-1]
	}

	return input[index]
}
