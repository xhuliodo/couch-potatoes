package application_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
)

func TestRateMovieRejectsInvalidRatingWithoutWriting(t *testing.T) {
	for _, rating := range []int{-1, 2, 100} {
		repo := db.NewInMemoryRepository()
		ratingService := application.NewRatingService(repo)

		err := ratingService.RateMovie("10", "6", rating)
		if err == nil {
			t.Errorf("rating %d: want a client error, got nil", rating)
		} else if cause := errors.Cause(err); cause.Error() != "bad_request" {
			t.Errorf("rating %d: want cause bad_request, got %s", rating, cause)
		}

		user, err := repo.GetUserById("10")
		if err != nil {
			t.Fatalf("rating %d: could not get user: %s", rating, err)
		}
		if len(user.RatedMovies) != 0 {
			t.Errorf("rating %d: want no rating persisted, got %+v", rating, user.RatedMovies)
		}
	}
}

func TestRateMoviePersistsLikedAndNotLikedRatings(t *testing.T) {
	repo := db.NewInMemoryRepository()
	ratingService := application.NewRatingService(repo)

	if err := ratingService.RateMovie("10", "6", 1); err != nil {
		t.Fatalf("could not rate a movie as liked: %s", err)
	}
	if err := ratingService.RateMovie("10", "4", 0); err != nil {
		t.Fatalf("could not rate a movie as not liked: %s", err)
	}

	liked, err := storedRating(repo, "10", "6")
	if err != nil {
		t.Fatalf("could not read liked rating: %s", err)
	}
	if liked != 1 {
		t.Errorf("want stored rating 1 for liked movie, got %v", liked)
	}

	notLiked, err := storedRating(repo, "10", "4")
	if err != nil {
		t.Fatalf("could not read not-liked rating: %s", err)
	}
	if notLiked != 0 {
		t.Errorf("want stored rating 0 for not-liked movie, got %v", notLiked)
	}

	if err := ratingService.RateMovie("10", "6", 0); err != nil {
		t.Fatalf("could not rate the same movie again: %s", err)
	}

	user, err := repo.GetUserById("10")
	if err != nil {
		t.Fatalf("could not get user after re-rating: %s", err)
	}
	if len(user.RatedMovies) != 2 {
		t.Errorf("want a re-rating to update, not duplicate, the stored rating, got %d ratings", len(user.RatedMovies))
	}
	updated, err := storedRating(repo, "10", "6")
	if err != nil {
		t.Fatalf("could not read updated rating: %s", err)
	}
	if updated != 0 {
		t.Errorf("want updated stored rating 0, got %v", updated)
	}
}

func storedRating(repo *db.InMemoryRepository, userId, movieId string) (float64, error) {
	user, err := repo.GetUserById(userId)
	if err != nil {
		return 0, err
	}
	for _, ratedMovie := range user.RatedMovies {
		if ratedMovie.Id == movieId {
			return ratedMovie.Rating, nil
		}
	}
	return 0, errors.Errorf("movie %s has no stored rating", movieId)
}
