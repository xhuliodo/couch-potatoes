package domain

import (
	"errors"

	"github.com/google/uuid"
)

type MovieWatcher struct {
	Id             uuid.UUID
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

func (mw *MovieWatcher) GivesGenrePreferences(g []Genre) error {
	if len(g) < minimumGenreRequired {
		return errors.New("You have to select at least 3 genres to continue")
	}
	mw.FavoriteGenres = append(mw.FavoriteGenres, g...)
	return nil
}

func (mw *MovieWatcher) GetRatingHistory() ([]RatedMovie, error) {
	if len(mw.RatedMovies) == 0 {
		return []RatedMovie{}, errors.New("You have not rated any movies yet")
	}
	return mw.RatedMovies, nil
}
