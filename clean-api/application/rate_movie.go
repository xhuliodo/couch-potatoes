package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type RatingService struct {
	movieRepo domain.MovieRepo
	userRepo  UserRepo
}

func NewRatingService(movieRepo domain.MovieRepo, userRepo UserRepo) RatingService {
	return RatingService{movieRepo, userRepo}
}

func (rms RatingService) RateMovie(userId, movieId string, rating int) error {
	if _, err := rms.userRepo.GetUserById(userId); err != nil {
		return errors.New("a user with this identifier does not exist")
	}

	if _, err := rms.movieRepo.GetMovieById(movieId); err != nil {
		return errors.New("a movie with this identifier does not exist")
	}

	if err := rms.userRepo.RateMovie(userId, movieId, rating); err != nil {
		return errors.New("the rating was not successful, try again")
	}

	return nil
}
