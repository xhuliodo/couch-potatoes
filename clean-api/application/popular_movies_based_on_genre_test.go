package application_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func TestSorting(t *testing.T) {

	first := []float32{.5, .5, .5, .5, .5, .5, .5, .5, .5, .5}
	second := []float32{.5, .5, .5, .5, .5, .5}
	third := []float32{.1, .9, .5}
	fourth := []float32{1., .1, .1}
	unsortedList := []domain.AggregateMovieRatings{
		{
			Movie:      domain.Movie{},
			AllRatings: second,
		}, {
			Movie:      domain.Movie{},
			AllRatings: first,
		}, {
			Movie:      domain.Movie{},
			AllRatings: third,
		},
		{
			Movie:      domain.Movie{},
			AllRatings: fourth,
		},
	}

	want := []domain.PopulatiryScoredMovie{
		{
			Movie:        domain.Movie{},
			AvgRating:    .5,
			CountRatings: 10,
		}, {
			Movie:        domain.Movie{},
			AvgRating:    .5,
			CountRatings: 6,
		}, {
			Movie:        domain.Movie{},
			AvgRating:    .5,
			CountRatings: 3,
		}, {
			Movie:        domain.Movie{},
			AvgRating:    .4,
			CountRatings: 3,
		},
	}

	sortedList := application.SortMoviesBasedOnRatings(unsortedList)

	if !reflect.DeepEqual(sortedList, want) {
		t.Errorf("got %v, want %v", sortedList, want)
	}

}
