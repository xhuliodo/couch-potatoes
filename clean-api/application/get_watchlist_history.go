package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func (gws GetWatchlistService) GetWatchlistHistory(userId string, limit uint, skip uint) ([]domain.Watchlist, error) {
	emptyWatchlistHistory := []domain.Watchlist{}
	if _, err := gws.repo.GetUserById(userId); err != nil {
		return emptyWatchlistHistory, errors.New("a user with this identifier does not exist")
	}

	watchlistHistory, err := gws.repo.GetWatchlistHistory(userId)
	if err != nil {
		return emptyWatchlistHistory, errors.New("there are no more movies in your watchlist history")
	}

	watchlistHistory.SortByTimeAdded()

	// handle pagination
	length := len(watchlistHistory)
	begin, end, err := handlePagination(uint(length), skip, limit)
	if err != nil {
		return emptyWatchlistHistory, errors.New("you're all caught up")
	}
	userWatchlistHistory := watchlistHistory[begin:end]

	moviesIds := getWatchlistMovieIds(userWatchlistHistory)
	moviesDetails, _ := gws.repo.GetMoviesDetails(moviesIds)

	userWatchlistHistory.PopulateMoviesWithDetails(moviesDetails)

	return userWatchlistHistory, nil
}
