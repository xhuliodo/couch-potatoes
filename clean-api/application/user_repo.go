package application

import (
	"github.com/google/uuid"
)

type UserRepo interface {
	SaveGenrePreferences(userId string, genresId []uuid.UUID) error
	GetUserById(userId string) (User, error)
}
