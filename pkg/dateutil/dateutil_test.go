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
		numberOfDaysBetween int
	}{
		{asDate(2022, time.December, 4), asDate(2022, time.December, 9), 5},   // monday to friday
		{asDate(2022, time.December, 4), asDate(2022, time.December, 12), 6},  // monday to monday
		{asDate(2022, time.December, 12), asDate(2022, time.December, 4), 6},  // monday to monday (reversed)
		{asDate(2022, time.December, 13), asDate(2022, time.December, 21), 7}, // tuesday to wednesday
		{asDate(2022, time.December, 4), asDate(2022, time.December, 4), 0},   // monday to monday
		{asDate(2022, time.December, 9), asDate(2022, time.December, 12), 2},  // friday to monday
	}

	for _, table := range tables {
		result := WeekDaysBetween(table.inputDate1, table.inputDate2)
		if result != table.numberOfDaysBetween {
			t.Errorf("Weekdays between %s and %s calculated as %d, expected %d", table.inputDate1, table.inputDate2, result, table.numberOfDaysBetween)
		}
	}
}

func asDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
