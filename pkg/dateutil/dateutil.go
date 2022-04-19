package dateutil

import (
	"math"
	"time"
)

func DaysBetween(time1 time.Time, time2 time.Time) int {
	t1, t2 := time1, time2
	if time1.After(time2) {
		t2, t1 = time1, time2
	}
	return int(math.Ceil(t2.Sub(t1).Hours() / 24))
}
