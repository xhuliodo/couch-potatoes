package domain

import "github.com/pkg/errors"

var ErrNotFound = errors.New("No genre found :`(")

type MovieRepo interface {
	GetAllGenres() ([]Genre, error)
}
