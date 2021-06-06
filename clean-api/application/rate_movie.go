package application

import (
	"github.com/pkg/errors"
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
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return errStack
	}

	if _, err := rms.repo.GetMovieById(movieId); err != nil {
		errStack := errors.Wrap(err, "a movie with this identifier does not exist")
		return errStack
	}

	if err := rms.repo.RateMovie(userId, movieId, rating); err != nil {
		errStack := errors.Wrap(err, "rating a movie was not successful")
		return errStack
	}

	return nil
}
