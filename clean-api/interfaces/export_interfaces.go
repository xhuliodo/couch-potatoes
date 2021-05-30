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

func NewPopularMoviesInterface(repo domain.Repository) popularMoviesResource {
	popularMoviesService := application.NewPopularMovieService(repo)
	return popularMoviesResource{popularMoviesService}
}

func NewUserBasedRecInterface(repo domain.Repository) userBasedRecResource {
	userBasedRecService := application.NewUserBasedRecommendationService(repo)
	return userBasedRecResource{userBasedRecService}
}

func NewContentBasedRecInterface(repo domain.Repository) contentBasedRecResource {
	contentBasedRecService := application.NewContentBasedRecommendationService(repo)
	return contentBasedRecResource{contentBasedRecService}
}

func NewAddToWatchlistInterface(repo domain.Repository) addToWatchlistResource {
	addToWatchlistService := application.NewAddToWatchlistService(repo)
	return addToWatchlistResource{addToWatchlistService}
}

func NewRemoveFromWatchlistInterface(repo domain.Repository) removeFromWatchlistResource {
	removeFromWatchlistService := application.NewRemoveFromWatchlistService(repo)
	return removeFromWatchlistResource{removeFromWatchlistService}
}
