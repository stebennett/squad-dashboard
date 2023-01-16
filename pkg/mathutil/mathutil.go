package mathutil

import (
	"math"
	"sort"
)

type XY struct {
	X, Y float64
}

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

func LinearRegression(inData []XY) (r []XY, m float64, b float64) {
	q := len(inData)

	if q == 0 {
		return make([]XY, 0), 0, 0
	}

	p := float64(q)

	sum_x, sum_y, sum_xx, sum_xy := 0.0, 0.0, 0.0, 0.0

	for _, p := range inData {
		sum_x += p.X
		sum_y += p.Y
		sum_xx += p.X * p.X
		sum_xy += p.X * p.Y
	}

	m = (p*sum_xy - sum_x*sum_y) / (p*sum_xx - sum_x*sum_x)
	b = (sum_y / p) - (m * sum_x / p)

	r = make([]XY, q)
	for i, p := range inData {
		r[i].X = p.X
		r[i].Y = p.X*m + b
	}

	return r, m, b
}
