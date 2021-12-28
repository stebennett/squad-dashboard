package paginator

import "errors"

type PaginationArgs struct {
	startAt int
	perPage int
}

func nextPaginationArgs(currentStartAt int, totalPerPage int, currentCount int, totalCount int) (*PaginationArgs, error) {
	return nil, errors.New("Reached end of pagination")
}
