package domain

import "errors"

var ErrCouldNotSaveGenrePref = errors.New("Could not save genres preferences :`(")

type MovieWatcherRepo interface {
	SaveGenrePreferences(MovieWatcherID, []Genre) error
}