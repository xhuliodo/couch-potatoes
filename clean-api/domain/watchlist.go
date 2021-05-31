package domain

import "sort"

type UserWatchlist []Watchlist

type Watchlist struct {
	Movie
	TimeAdded int64
	Rating    float64
}

func (uw UserWatchlist) SortByTimeAdded() {
	sort.SliceStable(uw, func(i, j int) bool {
		return uw[i].TimeAdded > uw[j].TimeAdded
	})
}

func (uw UserWatchlist) PopulateMoviesWithDetails(moviesDetails MoviesDetails) {
	for i, movie := range uw {
		uw[i].Movie.Title = moviesDetails[movie.Id].Title
		uw[i].Movie.ReleaseYear = moviesDetails[movie.Id].ReleaseYear
		uw[i].Movie.MoreInfoLink = moviesDetails[movie.Id].MoreInfoLink
	}
}
