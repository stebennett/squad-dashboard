package mathutil

import (
	"math"
	"sort"
)

func Percentile(percentile float64, input []int) int {
	if len(input) == 0 {
		return 0
	}

	sort.Ints(input)
	index := int(math.Round(percentile * float64(len(input))))
	return input[index]
}
