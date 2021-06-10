package application_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
)

func TestFilterBasedOnRatingCount(t *testing.T) {

	repo := db.NewInMemoryRepository()
	contentBasedRecService := application.NewContentBasedRecommendationService(repo)

	userId := "7"
	got, _ := contentBasedRecService.GetContentBasedRecommendation(userId, 5)
	scores := []float64{}
	for _, g := range got {
		scores = append(scores, g.Score)
	}

	expectedScores := []float64{1, .5}

	if !reflect.DeepEqual(scores, expectedScores) {
		t.Errorf("got %v, want %v", scores, expectedScores)
	}

}
