package interfaces

import (
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func NewInitialSetupInterface(repo domain.Repository) setupResource {
	setupService := application.NewSetupService(repo)
	return setupResource{setupService}
}

func NewRateMovieInterface(repo domain.Repository) ratingResource {
	ratingService := application.NewRatingService(repo)
	return ratingResource{ratingService}
}

func NewRegisterUserInterface(repo domain.Repository) registerResource {
	registerService := application.NewRegisterService(repo)
	return registerResource{registerService}
}
