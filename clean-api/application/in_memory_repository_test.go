package application_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
)

func TestSaveGenrePreferencesPersistsOnUser(t *testing.T) {
	repo := db.NewInMemoryRepository()
	setupService := application.NewSetupService(repo)

	before, err := setupService.GetSetupStep("10")
	if err != nil {
		t.Fatalf("could not get setup step before saving genres: %s", err)
	}
	if before.Step != 1 {
		t.Fatalf("want user to start at step 1, got %d", before.Step)
	}

	if err := setupService.SaveGenrePreferences("10", []string{"1", "2", "3"}); err != nil {
		t.Fatalf("could not save genre preferences: %s", err)
	}

	after, err := setupService.GetSetupStep("10")
	if err != nil {
		t.Fatalf("could not get setup step after saving genres: %s", err)
	}
	if after.Step != 2 || after.RatingsGiven != 0 {
		t.Errorf("want user to move to step 2 with no ratings, got %+v", after)
	}
}

func TestGetMovieByIdReturnsStoredMovieOrNotFound(t *testing.T) {
	repo := db.NewInMemoryRepository()

	got, err := repo.GetMovieById("1")
	if err != nil {
		t.Fatalf("could not get known movie: %s", err)
	}
	want := domain.Movie{
		Title:          "Toy Story",
		ReleaseYear:    2000,
		MoreInfoLink:   "nana",
		PeopleInvolved: []domain.Cast{{Id: "2"}, {Id: "3"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}

	if _, err := repo.GetMovieById("999"); err == nil {
		t.Error("want an error for an unknown movie id, got nil")
	}
}
