package domain_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func TestCalculateJaccard(t *testing.T) {
	gotUsersLikedMovie := make(domain.UsersLikedMovies)
	populateWithStubUsersLikedMovies(&gotUsersLikedMovie)

	gotSimilarMoviesToLikedOnes := make(domain.SimilarMoviesToLikedOnes)
	populateWithStubSimilarMoviesToLikedOnes(&gotSimilarMoviesToLikedOnes)

	jacRes, _ := domain.CalculateJaccard(gotUsersLikedMovie, gotSimilarMoviesToLikedOnes)
	jacRes.SortByScoreDesc()
	jacResWithNoDupes := jacRes.RemoveDuplicates()

	expectedRes := domain.ContentBasedRecommendations{
		{domain.Movie{Id: "8"}, 1},
		{domain.Movie{Id: "7"}, 0.4},
		{domain.Movie{Id: "6"}, 0.2857142857142857},
	}

	if !reflect.DeepEqual(jacResWithNoDupes, expectedRes) {
		t.Errorf("\n got %+v \n want %+v", jacResWithNoDupes, expectedRes)
	}
}

func populateWithStubUsersLikedMovies(ulm *domain.UsersLikedMovies) {
	castDetails := [][]string{
		{"10", "12", "13", "14", "15"},
		{"12", "15", "18", "19", "20"},
	}
	movieIds := []string{"1", "2"}
	for _, movieId := range movieIds {
		(*ulm)[movieId] = domain.UsersLikedMovie{AllCast: map[string]bool{}}
	}
	for i, movieId := range movieIds {
		(*ulm).PopulateWithCast(movieId, castDetails[i])
	}
}

func populateWithStubSimilarMoviesToLikedOnes(smlo *domain.SimilarMoviesToLikedOnes) {
	castDetails := [][]string{
		{"10", "12", "21", "22"},
		{"12", "15"},
		{"10", "12", "13", "14", "15"},
	}
	movieIds := []string{"6", "7", "8"}
	for _, movieId := range movieIds {
		(*smlo)[movieId] = domain.SimilarMovieToLikedOnes{}
	}
	for i, movieId := range movieIds {
		(*smlo).PopulateWithCast(movieId, castDetails[i])
	}
}
