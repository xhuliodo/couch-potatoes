package interfaces

import (
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func NewInitialSetupInterface(movieRepo domain.MovieRepo, userRepo application.UserRepo) setupResource {
	setupService := application.NewSetupService(movieRepo, userRepo)
	return setupResource{setupService}
}

func NewRateMovieInterface(movieRepo domain.MovieRepo, userRepo application.UserRepo) ratingResource {
	ratingService := application.NewRatingService(movieRepo, userRepo)
	return ratingResource{ratingService}
}

func NewRegisterUserInterface(userRepo application.UserRepo) registerResource {
	registerService := application.NewRegisterService(userRepo)
	return registerResource{registerService}
}
