package domain

import (
	"errors"
)

type User struct {
	Id             string
	Name           string
	IsAdmin        bool
	FavoriteGenres []Genre
	RatedMovies    []RatedMovie
	Watchlist      []Movie
}

type RatedMovie struct {
	Movie
	Rating float64
}

const minimumGenreRequired int = 3

func (mw *User) GiveGenrePreferences(g []Genre) error {
	if len(g) < minimumGenreRequired {
		return errors.New("you have to select at least 3 genres to continue")
	}
	return nil
}
