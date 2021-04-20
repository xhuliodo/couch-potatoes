package use_cases

import "github.com/xhuliodo/couch-potatoes/clean-api/domain"

type User struct {
	isAdmin      bool
	MovieWatcher domain.MovieWatcher
}

