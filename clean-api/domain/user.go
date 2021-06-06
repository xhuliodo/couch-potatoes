package domain

import "github.com/pkg/errors"

type User struct {
	Id             string
	Name           string
	IsAdmin        bool
	FavoriteGenres []Genre
	RatedMovies    []RatedMovie
	Watchlists     []Watchlist
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
