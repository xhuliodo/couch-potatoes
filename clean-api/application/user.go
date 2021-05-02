package application

import (
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type User struct {
	Id           string
	IsAdmin      bool
	MovieWatcher domain.MovieWatcher
}

