package application

import (
	"errors"

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
		return errors.New("a user with this identifier does not exist")
	}

	if _, err := rfws.repo.GetMovieById(movieId); err != nil {
		return errors.New("a movie with this identifier does not exist")
	}

	if err := rfws.repo.RemoveFromWatchlist(userId, movieId); err != nil {
		return errors.New("movie was not removed from watchlist, try again")
	}

	return nil
}
