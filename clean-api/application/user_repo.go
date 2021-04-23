package application

import (
	"github.com/google/uuid"
)

type UserRepo interface {
	SaveGenrePreferences(userId UserId, genresId []uuid.UUID) error
	ById(userId UserId) (User, error)
}
