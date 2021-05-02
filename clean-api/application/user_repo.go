package application

import "github.com/xhuliodo/couch-potatoes/clean-api/domain"

type UserRepo interface {
	SaveGenrePreferences(userId string, genres []domain.Genre) error
	GetUserById(userId string) (User, error)
}
