package application

import (
	"time"

	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type AddToWatchlistService struct {
	repo domain.Repository
}

func NewAddToWatchlistService(repo domain.Repository) AddToWatchlistService {
	return AddToWatchlistService{repo}
}

func (atws AddToWatchlistService) AddToWatchlist(userId, movieId string) error {
	if _, err := atws.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return errStack
	}

	if _, err := atws.repo.GetMovieById(movieId); err != nil {
		errStack := errors.Wrap(err, "a movie with this identifier does not exist")
		return errStack
	}

	now := time.Now()
	timeOfAdding := now.Unix()
	if err := atws.repo.AddToWatchlist(userId, movieId, timeOfAdding); err != nil {
		errStack := errors.Wrap(err, "movie was not added to watchlist")
		return errStack
	}

	return nil
}
