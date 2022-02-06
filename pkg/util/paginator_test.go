package util

import "testing"

func TestNextPaginationArgs(t *testing.T) {

	tables := []struct {
		currentStartAt int
		totalPerPage   int
		currentCount   int
		totalCount     int
		newStartAt     int
	}{
		{0, 100, 100, 1000, 100},
		{100, 100, 100, 100, -1},
		{0, 50, 25, 25, -1},
		{0, 20, 0, 0, -1},
	}

	for _, table := range tables {
		result := NextPaginationArgs(table.currentStartAt, table.totalPerPage, table.currentCount, table.totalCount)

		if result != table.newStartAt {
			t.Errorf("New StartAt value was incorrect, got %d, expected %d", result, table.newStartAt)
		}
	}
}
