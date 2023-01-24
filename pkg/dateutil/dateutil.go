package dateutil

import (
	"fmt"
	"math"
	"time"
)

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

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

func PreviousWeekDates(startDate time.Time, numberOfWeeks int) []time.Time {
	if numberOfWeeks == 0 {
		return []time.Time{}
	}

	outDates := []time.Time{startDate}
	for i := 1; i < numberOfWeeks; i++ {
		outDates = append(outDates, startDate.AddDate(0, 0, i*-7))
	}
	return outDates
}

func AsDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func ParseDayOfWeek(s string) (time.Weekday, error) {
	if d, ok := daysOfWeek[s]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", s)
}
