package application_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func TestCalculateScore(t *testing.T) {
	calculateScore := application.CalculateScore

	usersToCompare := make(domain.UsersToCompare)
	populateWithStubUsersToCompare(&usersToCompare)
	moviesToScore := make(domain.ScoringMovies)
	populateWithStubMoviesToScore(&moviesToScore)

	got := calculateScore(usersToCompare, moviesToScore)

	sort.SliceStable(got, func(i, j int) bool {
		return got[i].Score > got[j].Score
	})

	want := domain.UsersBasedRecommendation{
		{Movie: domain.Movie{Id: "10"}, Score: 0.5},
		{Movie: domain.Movie{Id: "12"}, Score: 0.215},
		{Movie: domain.Movie{Id: "11"}, Score: 0.2},
		{Movie: domain.Movie{Id: "13"}, Score: 0.1},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n got %v,\n want: %v", got, want)
	}
}

func populateWithStubUsersToCompare(utc *domain.UsersToCompare) {
	usersId := []string{"1", "2", "3"}
	pearsons := []float64{1, 0.5, 0.2}
	for i, userId := range usersId {
		userDetails := domain.UserToCompare{PearsonCoefficient: pearsons[i]}
		(*utc)[userId] = &userDetails
	}
}

func populateWithStubMoviesToScore(sc *domain.ScoringMovies) {
	movieIds := []string{"10", "11", "12", "13"}
	ratings := []domain.Rating{
		map[string]float64{"1": 0.5, "2": 1},
		map[string]float64{"3": 1},
		map[string]float64{"2": 0.7, "3": 0.4},
		map[string]float64{"2": 0.2},
	}
	for i, movieId := range movieIds {
		movieDetail := domain.Details{Ratings: ratings[i]}
		(*sc)[movieId] = &movieDetail
	}

}
