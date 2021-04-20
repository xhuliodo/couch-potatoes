package domain

import "errors"

var ErrNotFound = errors.New("No genre found :`(")

type MovieRepo interface {
	GetAllGenres() ([]Genre, error)
}
