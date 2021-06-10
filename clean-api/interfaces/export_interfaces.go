package interfaces

import (
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func NewInitialSetupInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) setupResource {
	setupService := application.NewSetupService(repo)
	return setupResource{setupService}
}

func NewRateMovieInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) ratingResource {
	ratingService := application.NewRatingService(repo)
	return ratingResource{ratingService}
}

func NewRegisterUserInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) registerResource {
	registerService := application.NewRegisterService(repo)
	return registerResource{registerService}
}

func NewPopularMoviesInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) popularMoviesResource {
	popularMoviesService := application.NewPopularMovieService(repo)
	return popularMoviesResource{popularMoviesService}
}

func NewUserBasedRecInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) userBasedRecResource {
	userBasedRecService := application.NewUserBasedRecommendationService(repo)
	return userBasedRecResource{userBasedRecService}
}

func NewContentBasedRecInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) contentBasedRecResource {
	contentBasedRecService := application.NewContentBasedRecommendationService(repo)
	return contentBasedRecResource{contentBasedRecService}
}

func NewAddToWatchlistInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) addToWatchlistResource {
	addToWatchlistService := application.NewAddToWatchlistService(repo)
	return addToWatchlistResource{addToWatchlistService}
}

func NewRemoveFromWatchlistInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) removeFromWatchlistResource {
	removeFromWatchlistService := application.NewRemoveFromWatchlistService(repo)
	return removeFromWatchlistResource{removeFromWatchlistService}
}

func NewGetWatchlistInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,

) getWatchlistResource {
	getWatchlistService := application.NewGetWatchlistService(repo)
	return getWatchlistResource{getWatchlistService}
}

func NewHealthcheckInterface(
	repo domain.Repository,
	errorLogger domain.ErrorLoggerInterface,
) healthcheckResource {
	healthcheckService := application.NewHealthcheckService(repo)
	return healthcheckResource{healthcheckService, errorLogger}
}
