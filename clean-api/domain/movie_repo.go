package domain

import "errors"

var ErrNotFound = errors.New("no genre found :`(")

type MovieRepo interface {
	GetAllGenres() ([]Genre, error)
}
