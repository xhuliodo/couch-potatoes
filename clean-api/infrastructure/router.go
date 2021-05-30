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

	router.Get("/genres", initialSetupInterface.GetAllGenres)
	r := router.Route("/users", func(r chi.Router) {
		r.Post("/", registerUserInterface.RegisterUser)
		r.Get("/setup", initialSetupInterface.GetUserSetupStep)
		r.Post("/genres", initialSetupInterface.SaveGenrePreferences)
		r.Post("/ratings", rateMovieInterface.RateMovie)
		r.Post("/watchlist/{movieId}", addToWatchlistInterface.AddToWatchlist)
		r.Delete("/watchlist/{movieId}", removeFromWatchlistInterface.RemoveFromWatchlist)
		// todo: these two are with pagination
		// r.Get("/watchlist")
		// r.Get("/watchlist-history")
	})

	r.Route("/recommendations", func(r chi.Router) {
		router.Get("/recommendations/popular", popularMoviesInterface.GetPopularMoviesBasedOnGenre)
		router.Get("/recommendations/user-based", userBasedRecInterface.GetUserBasedRecommendation)
		router.Get("/recommendations/content-based", contentBasedRecInterface.GetContentBasedRecommendation)
	})

}
