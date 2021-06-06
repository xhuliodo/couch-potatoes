// Clean Potatoes API
//
// Docs for Movies API
//
// 	Schemes: http
// 	Host: localhost
// 	BasePath: /
// 	Version: 1.0.0
//
// 	Consumes:
// 	- application/json
//
// 	Produces:
// 	- application/json
// swagger:meta
package infrastructure

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xhuliodo/couch-potatoes/clean-api/interfaces"
)

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(authMiddleware)

	return r
}

func CreateRoutes(router *chi.Mux, repo *Neo4jRepository) {
	initialSetupInterface := interfaces.NewInitialSetupInterface(repo)
	rateMovieInterface := interfaces.NewRateMovieInterface(repo)
	registerUserInterface := interfaces.NewRegisterUserInterface(repo)
	popularMoviesInterface := interfaces.NewPopularMoviesInterface(repo)
	userBasedRecInterface := interfaces.NewUserBasedRecInterface(repo)
	contentBasedRecInterface := interfaces.NewContentBasedRecInterface(repo)
	addToWatchlistInterface := interfaces.NewAddToWatchlistInterface(repo)
	removeFromWatchlistInterface := interfaces.NewRemoveFromWatchlistInterface(repo)
	getWatchlistInterface := interfaces.NewGetWatchlistInterface(repo)

	// movie routes
	router.Get("/genres", initialSetupInterface.GetAllGenres)

	// user routes
	router.Route("/users", func(r chi.Router) {
		r.Post("/{userId}", registerUserInterface.RegisterUser)
		r.Get("/setup", initialSetupInterface.GetUserSetupStep)
		r.Post("/genres", initialSetupInterface.SaveGenrePreferences)
		r.Post("/ratings", rateMovieInterface.RateMovie)
	})

	// recommendation routes
	router.Route("/recommendations", func(r chi.Router) {
		r.Get("/popular", popularMoviesInterface.GetPopularMoviesBasedOnGenre)
		r.Get("/user-based", userBasedRecInterface.GetUserBasedRecommendation)
		r.Get("/content-based", contentBasedRecInterface.GetContentBasedRecommendation)
	})

	// watchlist routes
	router.Route("/watchlist", func(r chi.Router) {
		r.Get("/", getWatchlistInterface.GetWatchlist)
		r.Post("/{movieId}", addToWatchlistInterface.AddToWatchlist)
		r.Delete("/{movieId}", removeFromWatchlistInterface.RemoveFromWatchlist)
		r.Get("/history", getWatchlistInterface.GetWatchlistHistory)
	})

}
