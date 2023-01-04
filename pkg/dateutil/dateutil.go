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

func WeekDaysBetween(time1 time.Time, time2 time.Time, datesToExclude []time.Time) int {
	// check the ordering of the dates
	t1, t2 := time1, time2
	if time1.After(time2) {
		t2, t1 = time1, time2
	}

	// loop through the days and count number of non-weekend days
	dayCounter := 0
	incDate := t1
	for {
		if incDate.After(t2) {
			break
		}

		dayOfWeek := incDate.Weekday()
		if dayOfWeek != time.Saturday && dayOfWeek != time.Sunday && !ContainsDate(incDate, datesToExclude) {
			dayCounter++
		}

		incDate = incDate.AddDate(0, 0, 1)
	}

	return dayCounter
}

func NearestPreviousDateForDay(inTime time.Time, targetDay time.Weekday) time.Time {
	if inTime.Weekday() == targetDay {
		return inTime
	} else if inTime.Weekday() > targetDay {
		daysToPreviousTargetDay := int(inTime.Weekday()) - int(targetDay)
		return inTime.AddDate(0, 0, -1*daysToPreviousTargetDay)
	} else {
		daysToMoveToNextFriday := int(targetDay) - int(inTime.Weekday())
		return inTime.AddDate(0, 0, daysToMoveToNextFriday).AddDate(0, 0, -7)
	}
}

func ContainsDate(needle time.Time, haystack []time.Time) bool {
	for _, d := range haystack {
		if d.Year() == needle.Year() && d.Month() == needle.Month() && d.Day() == needle.Day() {
			return true
		}
	}
	return false
}
