package dateutil

import (
	"testing"
	"time"
)

func TestLastCompleteWeekFriday(t *testing.T) {
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
		result := LastCompleteWeekFriday(table.inputDate)
		if !result.Equal(table.outputDate) {
			t.Errorf("Last Friday Date incorrect for input %s, got %s, expected %s", table.inputDate, result.String(), table.outputDate.String())
		}
	}

}

func asDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
