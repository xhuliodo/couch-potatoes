package domain_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func TestFilterBasedOnRatingCount(t *testing.T) {
	var r float64 = 0
	got := domain.UsersToCompare{
		"1": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
		"2": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
		"3": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
	}

	want := domain.UsersToCompare{
		"1": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
		"3": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
	}

	got.FilterBasedOnRatingsCount()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}
