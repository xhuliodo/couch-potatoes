package domain_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func TestCalculateBoostedScore(t *testing.T) {
	popularMovies := domain.PopularMovies{
		{RatingsCount: 10, GenreMatched: 1},
		{RatingsCount: 20, GenreMatched: 5},
		{RatingsCount: 10, GenreMatched: 10},
		{RatingsCount: 20, GenreMatched: 1},
	}

	popularMovies.CalculateBoostedScore()

	expectedScores := domain.PopularMovies{
		{RatingsCount: 10, GenreMatched: 1, BoostedScore: 10},
		{RatingsCount: 20, GenreMatched: 5, BoostedScore: 100},
		{RatingsCount: 10, GenreMatched: 10, BoostedScore: 100},
		{RatingsCount: 20, GenreMatched: 1, BoostedScore: 20},
	}

	if !reflect.DeepEqual(popularMovies, expectedScores) {
		t.Errorf("got: %+v, wanted: %+v", popularMovies, expectedScores)
	}
}
