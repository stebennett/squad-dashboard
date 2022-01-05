package util

func nextPaginationArgs(currentStartAt int, totalPerPage int, currentCount int, totalCount int) int {
	if currentStartAt + currentCount < totalCount {
		return currentStartAt + totalPerPage
	}
	return -1
}
