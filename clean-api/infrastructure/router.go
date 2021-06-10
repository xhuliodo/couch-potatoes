package infrastructure

import (
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/xhuliodo/couch-potatoes/clean-api/docs"

	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/auth"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/db"
	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/logger"
	"github.com/xhuliodo/couch-potatoes/clean-api/interfaces"
)

func CreateRouter(accessLogger *logrus.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.NewAccessLoggerMiddleware(accessLogger))
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	return r
}

func CreateRoutes(router *chi.Mux, repo *db.Neo4jRepository, errorLogger *logger.ErrorLogger) {
	initialSetupInterface := interfaces.NewInitialSetupInterface(repo, errorLogger)
	rateMovieInterface := interfaces.NewRateMovieInterface(repo, errorLogger)
	registerUserInterface := interfaces.NewRegisterUserInterface(repo, errorLogger)
	popularMoviesInterface := interfaces.NewPopularMoviesInterface(repo, errorLogger)
	userBasedRecInterface := interfaces.NewUserBasedRecInterface(repo, errorLogger)
	contentBasedRecInterface := interfaces.NewContentBasedRecInterface(repo, errorLogger)
	addToWatchlistInterface := interfaces.NewAddToWatchlistInterface(repo, errorLogger)
	removeFromWatchlistInterface := interfaces.NewRemoveFromWatchlistInterface(repo, errorLogger)
	getWatchlistInterface := interfaces.NewGetWatchlistInterface(repo, errorLogger)
	healthcheckInterface := interfaces.NewHealthcheckInterface(repo, errorLogger)

	// movie routes
	router.Route("/genres", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", initialSetupInterface.GetAllGenres)
	})

	// user routes
	router.Route("/users", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/{userId}", registerUserInterface.RegisterUser)
		r.Get("/setup", initialSetupInterface.GetUserSetupStep)
		r.Post("/genres", initialSetupInterface.SaveGenrePreferences)
		r.Post("/ratings", rateMovieInterface.RateMovie)
	})

	// recommendation routes
	router.Route("/recommendations", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/popular", popularMoviesInterface.GetPopularMoviesBasedOnGenre)
		r.Get("/user-based", userBasedRecInterface.GetUserBasedRecommendation)
		r.Get("/content-based", contentBasedRecInterface.GetContentBasedRecommendation)
	})

	// watchlist routes
	router.Route("/watchlist", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", getWatchlistInterface.GetWatchlist)
		r.Post("/{movieId}", addToWatchlistInterface.AddToWatchlist)
		r.Delete("/{movieId}", removeFromWatchlistInterface.RemoveFromWatchlist)
		r.Get("/history", getWatchlistInterface.GetWatchlistHistory)
	})

	router.Mount("/docs", httpSwagger.WrapHandler)

	router.Get("/healthcheck", healthcheckInterface.GetHealthcheck)
	router.Get("/ready", healthcheckInterface.GetReady)
}
