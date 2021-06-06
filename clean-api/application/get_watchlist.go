package application

import (
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type GetWatchlistService struct {
	repo domain.Repository
}

func NewGetWatchlistService(repo domain.Repository) GetWatchlistService {
	return GetWatchlistService{repo}
}

func (gws GetWatchlistService) GetWatchlist(userId string, limit uint, skip uint) (
	watchlist domain.UserWatchlist, err error,
) {
	if _, err := gws.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return watchlist, errStack
	}

	watchlist, err = gws.repo.GetWatchlist(userId)
	if err != nil {
		return watchlist, errors.Wrap(err, "could not get user watchlist")
	}

	if len(watchlist) < 1 {
		cause := errors.New("not_found")
		return watchlist, errors.Wrap(cause, "there are no more movies in your watchlist")
	}

	watchlist.SortByTimeAdded()

	// handle pagination
	length := len(watchlist)
	begin, end, err := handlePagination(uint(length), skip, limit)
	if err != nil {
		return watchlist, err
	}
	watchlist = watchlist[begin:end]

	moviesIds := getWatchlistMovieIds(watchlist)
	moviesDetails, _ := gws.repo.GetMoviesDetails(moviesIds)
	watchlist.PopulateMoviesWithDetails(moviesDetails)

	return watchlist, nil
}

func getWatchlistMovieIds(watchlist domain.UserWatchlist) []string {
	moviesIds := []string{}
	for _, w := range watchlist {
		movieId := w.Id
		moviesIds = append(moviesIds, movieId)
	}
	return moviesIds
}
