package application

import "github.com/pkg/errors"

const (
	defaultLimit uint = 5
	defaultSkip  uint = 0
)

func handlePagination(len, skip, limit uint) (begin, end uint, err error) {
	if skip != defaultSkip {
		maxSkip := maxSkip(uint(len), limit)
		if skip > maxSkip {
			cause := errors.New("bad_request")
			return 0, 0, errors.Wrap(cause, "you've all caught up, can't skip anymore than this")
		}
	}

	begin = skip

	remaining := len - skip
	if remaining < limit {
		end = remaining + skip
		return begin, end, nil
	}

	end = limit + skip
	return begin, end, nil

}

func maxSkip(total uint, limit uint) uint {
	if limit == 0 {
		limit = defaultLimit
	}
	return ((total - 1) / limit) * limit
}
