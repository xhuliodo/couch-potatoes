package application

import (
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type RemoveFromWatchlistService struct {
	repo domain.Repository
}

func NewRemoveFromWatchlistService(repo domain.Repository) RemoveFromWatchlistService {
	return RemoveFromWatchlistService{repo}
}

func (rfws RemoveFromWatchlistService) RemoveFromWatchlist(userId, movieId string) error {
	if _, err := rfws.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return errStack
	}

	if _, err := rfws.repo.GetMovieById(movieId); err != nil {
		errStack := errors.Wrap(err, "a movie with this identifier does not exist")
		return errStack
	}

	if err := rfws.repo.RemoveFromWatchlist(userId, movieId); err != nil {
		errStack := errors.Wrap(err, "movie was not removed from watchlist, try again")
		return errStack
	}

	return nil
}
