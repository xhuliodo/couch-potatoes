package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type GetWatchlistService struct {
	repo domain.Repository
}

func NewGetWatchlistService(repo domain.Repository) GetWatchlistService {
	return GetWatchlistService{repo}
}

func (gws GetWatchlistService) GetWatchlist(userId string, limit uint, skip uint) ([]domain.Watchlist, error) {
	emptyWatchlist := []domain.Watchlist{}
	if _, err := gws.repo.GetUserById(userId); err != nil {
		return emptyWatchlist, errors.New("a user with this identifier does not exist")
	}

	watchlist, err := gws.repo.GetWatchlist(userId)
	if err != nil {
		return emptyWatchlist, errors.New("there are no more movies in your watchlist")
	}

	watchlist.SortByTimeAdded()

	// handle pagination
	length := len(watchlist)
	begin, end, err := handlePagination(uint(length), skip, limit)
	if err != nil {
		return emptyWatchlist, errors.New("you're all caught up")
	}
	userWatchlist := watchlist[begin:end]

	moviesIds := getWatchlistMovieIds(userWatchlist)
	moviesDetails, _ := gws.repo.GetMoviesDetails(moviesIds)

	userWatchlist.PopulateMoviesWithDetails(moviesDetails)

	return userWatchlist, nil
}

func getWatchlistMovieIds(watchlist domain.UserWatchlist) []string {
	moviesIds := []string{}
	for _, w := range watchlist {
		movieId := w.Id
		moviesIds = append(moviesIds, movieId)
	}
	return moviesIds
}
