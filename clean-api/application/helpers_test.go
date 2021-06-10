package application_test

import (
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/application"
)

type paginationInput struct {
	len   uint
	skip  uint
	limit uint
}

type paginationResult struct {
	begin uint
	end   uint
}

func TestPaginationHelpers(t *testing.T) {
	inputs := []paginationInput{
		{15, 0, 0},
		{15, application.DefaultSkip, application.DefaultLimit},
		{15, 10, 20},
		{15, 15, 5},
		{15, 20, 5},
	}
	res := []paginationResult{
		{0, 0},
		{0, 5},
		{10, 15},
		{15, 15},
		{15, 15},
	}
	expectedErrString := "you've all caught up, can't skip anymore than this: bad_request"
	for i, input := range inputs {
		begin, end, err := application.HandlePagination(input.len, input.skip, input.limit)
		if begin != res[i].begin || end != res[i].end {
			gotErrString := err.Error()
			if expectedErrString!=gotErrString{
				t.Errorf(`with len %d, skip %d, limit %d 
			got begin %d, end %d
			want begin %d, end %d
			and err %s`,
				input.len, input.skip, input.limit,
				begin, end,
				res[i].begin, res[i].end,
				gotErrString)
			}
			
		}
	}

}
