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
			Movie:        domain.Movie{},
			AllRatings:   second,
			GenreMatched: 2,
		}, {
			Movie:        domain.Movie{},
			AllRatings:   first,
			GenreMatched: 1,
		}, {
			Movie:        domain.Movie{},
			AllRatings:   third,
			GenreMatched: 1,
		},
		{
			Movie:        domain.Movie{},
			AllRatings:   fourth,
			GenreMatched: 1,
		},
	}

	want := []domain.PopulatiryScoredMovie{
		{
			Movie:        domain.Movie{},
			AvgRating:    .5,
			CountRatings: 12,
		}, {
			Movie:        domain.Movie{},
			AvgRating:    .5,
			CountRatings: 10,
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

	sortedList, _ := application.SortMoviesBasedOnRatings(unsortedList)

	if !reflect.DeepEqual(sortedList, want) {
		t.Errorf("got %v, want %v", sortedList, want)
	}

}
