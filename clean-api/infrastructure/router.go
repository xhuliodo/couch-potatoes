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

	router.Get("/genres", initialSetupInterface.GetAllGenres)
	router.Route("/users", func(r chi.Router) {
		r.Post("/", registerUserInterface.RegisterUser)
		// r.Get("/{userId}", initialSetupInterface.GetUserById)
		r.Post("/genres", initialSetupInterface.SaveGenrePreferences)
		// TODO: get what setup step a user is at
		// r.Get("/{userId}/steps", ...)
		r.Post("/ratings", rateMovieInterface.RateMovie)
	})

	// TODO: figure out a way, if possible to group routes
	// router.Route("/recommendations", func(r chi.Router) {
	router.Get("/recommendations/popular", popularMoviesInterface.GetPopularMoviesBasedOnGenre)
	// router.Get("/user-based", )
	// router.Get("/content-based", )

	// })

}
