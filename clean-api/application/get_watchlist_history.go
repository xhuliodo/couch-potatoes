package application

import (
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func (gws GetWatchlistService) GetWatchlistHistory(userId string, limit uint, skip uint) (
	watchlistHistory domain.UserWatchlist, err error,
) {
	if _, err := gws.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return watchlistHistory, errStack
	}

	watchlistHistory, err = gws.repo.GetWatchlistHistory(userId)
	if err != nil {
		return watchlistHistory, errors.Wrap(err, "could not get user watchlist history")
	}

	if len(watchlistHistory) < 1 {
		cause := errors.New("not_found")
		return watchlistHistory, errors.Wrap(cause, "there are no movies in your watchlist history")
	}

	watchlistHistory.SortByTimeAdded()

	// handle pagination
	length := len(watchlistHistory)
	begin, end, err := handlePagination(uint(length), skip, limit)
	if err != nil {
		return watchlistHistory, err
	}
	watchlistHistory = watchlistHistory[begin:end]

	moviesIds := getWatchlistMovieIds(watchlistHistory)
	moviesDetails, _ := gws.repo.GetMoviesDetails(moviesIds)
	watchlistHistory.PopulateMoviesWithDetails(moviesDetails)

	return watchlistHistory, nil
}
