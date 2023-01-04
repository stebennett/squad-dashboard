package dateutil

import (
	"testing"
	"time"
)

func TestNearestPreviousDateForDay(t *testing.T) {
	tables := []struct {
		inputDate  time.Time
		outputDate time.Time
	}{
		{asDate(2022, time.July, 24), asDate(2022, time.July, 22)}, // sunday
		{asDate(2022, time.July, 25), asDate(2022, time.July, 22)}, // monday
		{asDate(2022, time.July, 26), asDate(2022, time.July, 22)}, // tuesday
		{asDate(2022, time.July, 27), asDate(2022, time.July, 22)}, // wednesday
		{asDate(2022, time.July, 28), asDate(2022, time.July, 22)}, // thursday
		{asDate(2022, time.July, 29), asDate(2022, time.July, 29)}, // friday
		{asDate(2022, time.July, 30), asDate(2022, time.July, 29)}, // saturday
	}

	for _, table := range tables {
		result := NearestPreviousDateForDay(table.inputDate, time.Friday)
		if !result.Equal(table.outputDate) {
			t.Errorf("Last Friday Date incorrect for input %s, got %s, expected %s", table.inputDate, result.String(), table.outputDate.String())
		}
	}

}

func TestWeekDaysBetween(t *testing.T) {
	tables := []struct {
		inputDate1          time.Time
		inputDate2          time.Time
		exclude             []time.Time
		numberOfDaysBetween int
	}{
		// simple cases - no excludes
		{asDate(2022, time.December, 4), asDate(2022, time.December, 9), []time.Time{}, 5},   // monday to friday
		{asDate(2022, time.December, 4), asDate(2022, time.December, 12), []time.Time{}, 6},  // monday to monday
		{asDate(2022, time.December, 12), asDate(2022, time.December, 4), []time.Time{}, 6},  // monday to monday (reversed)
		{asDate(2022, time.December, 13), asDate(2022, time.December, 21), []time.Time{}, 7}, // tuesday to wednesday
		{asDate(2022, time.December, 4), asDate(2022, time.December, 4), []time.Time{}, 0},   // monday to monday
		{asDate(2022, time.December, 9), asDate(2022, time.December, 12), []time.Time{}, 2},  // friday to monday
		// with excluded dates
		{asDate(2022, time.December, 4), asDate(2022, time.December, 12), []time.Time{asDate(2022, time.December, 6)}, 5},                                    // monday to monday
		{asDate(2022, time.December, 12), asDate(2022, time.December, 4), []time.Time{asDate(2022, time.December, 6)}, 5},                                    // monday to monday (reversed)
		{asDate(2022, time.December, 13), asDate(2022, time.December, 21), []time.Time{asDate(2022, time.December, 14), asDate(2022, time.December, 15)}, 5}, // tuesday to wednesday
	}

	for _, table := range tables {
		result := WeekDaysBetween(table.inputDate1, table.inputDate2, table.exclude)
		if result != table.numberOfDaysBetween {
			t.Errorf("Weekdays between %s and %s, excluding %s; calculated: %d, expected: %d", table.inputDate1, table.inputDate2, table.exclude, result, table.numberOfDaysBetween)
		}
	}
}

func TestContainsDate(t *testing.T) {
	tables := []struct {
		needle   time.Time
		haystack []time.Time
		result   bool
	}{
		{asDate(2022, 01, 01), []time.Time{}, false},
		{asDate(2022, 01, 01), []time.Time{asDate(2022, 01, 01)}, true},
		{asDate(2022, 01, 01), []time.Time{asDate(2022, 12, 01)}, false},
		{asDate(2022, 01, 01), []time.Time{asDate(2022, 12, 01), asDate(2022, 12, 07), asDate(2022, 12, 10)}, false},
		{asDate(2022, 12, 01), []time.Time{asDate(2022, 12, 01), asDate(2022, 12, 07), asDate(2022, 12, 10)}, true},
	}

	for _, table := range tables {
		result := ContainsDate(table.needle, table.haystack)
		if result != table.result {
			t.Errorf("Needle: %s; Haystack: %s; Expected: %t; Got: %t", table.needle, table.haystack, table.result, result)
		}
	}
}

func TestPreviousWeekDates(t *testing.T) {
	tables := []struct {
		startDate      time.Time
		numberOfWeeks  int
		expectedResult []time.Time
	}{
		{asDate(2022, 12, 30), 0, []time.Time{}},
		{asDate(2022, 12, 30), 1, []time.Time{asDate(2022, 12, 30)}},
		{asDate(2022, 12, 30), 3, []time.Time{asDate(2022, 12, 30), asDate(2022, 12, 23), asDate(2022, 12, 16)}},
	}

	for _, table := range tables {
		result := PreviousWeekDates(table.startDate, table.numberOfWeeks)
		if !containsExactly(table.expectedResult, result) {
			t.Errorf("Expected: %s; Got: %s", table.expectedResult, result)
		}
	}
}

func asDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func containsExactly(times1 []time.Time, times2 []time.Time) bool {
	if len(times1) != len(times2) {
		return false
	}

	for _, t := range times1 {
		if !ContainsDate(t, times2) {
			return false
		}
	}

	return true
}
