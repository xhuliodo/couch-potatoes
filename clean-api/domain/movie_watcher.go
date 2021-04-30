package domain

import (
	"errors"
)

type MovieWatcherID string

type MovieWatcher struct {
	Id             MovieWatcherID
	Name           string
	FavoriteGenres []Genre
	RatedMovies    []RatedMovie
	Watchlist      []Movie
}

type RatedMovie struct {
	Movie
	Rating float32
}

const minimumGenreRequired int = 3

func (mw *MovieWatcher) GiveGenrePreferences(g []Genre) error {
	if len(g) < minimumGenreRequired {
		return errors.New("you have to select at least 3 genres to continue")
	}
	mw.FavoriteGenres = append(mw.FavoriteGenres, g...)
	return nil
}
