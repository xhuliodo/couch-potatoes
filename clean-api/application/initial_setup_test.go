package application_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
)

type favoriteGenres []string

func TestSaveGenrePreferences(t *testing.T) {
	currentGenres := []domain.Genre{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
		{Id: "4"},
		{Id: "5"},
	}
	userFavoriteGenres := []favoriteGenres{
		{"1", "2", "3"},
		{"1", "5", "4"},
		{"1", "5", "12"},
	}

	want := [][]domain.Genre{
		{{Id: "1"}, {Id: "2"}, {Id: "3"}},
		{{Id: "1"}, {Id: "5"}, {Id: "4"}},
		{{Id: "1"}, {Id: "5"}},
	}

	for i, ufg := range userFavoriteGenres {
		got, err := application.AreGenreIdsValid(currentGenres, ufg)
		if !reflect.DeepEqual(got, want[i]) {
			t.Errorf("with currect genres %v, got %v, want %v, with err %s",
				currentGenres, got, want, err.Error())
		}
	}
}

func TestGetSetupStep(t *testing.T) {
	repo := db.NewInMemoryRepository()
	setupService := application.NewSetupService(repo)
	userIds := []string{"7", "10", "11"}
	wantedSetupStep := []application.SetupStep{
		{
			Step:         2,
			Finished:     false,
			Message:      "user has not given enough ratings to get to know him yet",
			RatingsGiven: 2,
		},
		{
			Step:     1,
			Finished: false,
			Message:  "user has yet to give genre preferences",
		},

		{
			Step:         3,
			Finished:     true,
			Message:      "user has finished the setup process",
			RatingsGiven: 16,
		},
	}
	for i, userId := range userIds {
		got, err := setupService.GetSetupStep(userId)
		if !reflect.DeepEqual(got, wantedSetupStep[i]) {
			errString := ""
			if err != nil {
				errString = err.Error()
			}
			t.Errorf("got %+v, want %+v, with err %s",
				got, wantedSetupStep[i], errString)
		}

	}
}
