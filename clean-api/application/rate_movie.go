package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type RatingService struct {
	repo domain.Repository
}

func NewRatingService(repo domain.Repository) RatingService {
	return RatingService{repo}
}

func (rms RatingService) RateMovie(userId, movieId string, rating int) error {
	if _, err := rms.repo.GetUserById(userId); err != nil {
		return errors.New("a user with this identifier does not exist")
	}

	if _, err := rms.repo.GetMovieById(movieId); err != nil {
		return errors.New("a movie with this identifier does not exist")
	}

	if err := rms.repo.RateMovie(userId, movieId, rating); err != nil {
		return errors.New("the rating was not successful, try again")
	}

	return nil
}
