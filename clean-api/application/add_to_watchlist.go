package application

import (
	"errors"
	"time"

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
		return errors.New("a user with this identifier does not exist")
	}

	if _, err := atws.repo.GetMovieById(movieId); err != nil {
		return errors.New("a movie with this identifier does not exist")
	}

	now := time.Now()
	timeOfAdding := now.Unix()
	if err := atws.repo.AddToWatchlist(userId, movieId, timeOfAdding); err != nil {
		return errors.New("movie was not added to watchlist, try again")
	}

	return nil
}
