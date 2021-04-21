package domain

import "github.com/pkg/errors"

var ErrCouldNotSaveGenrePref = errors.New("Could not save genres preferences :`(")

type MovieWatcherRepo interface {
	SaveGenrePreferences(MovieWatcherID, []Genre) error
}