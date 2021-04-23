package application

import (
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type UserId string

type User struct {
	Id           UserId
	isAdmin      bool
	MovieWatcher domain.MovieWatcher
}

