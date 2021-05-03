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
	initialSetupInterface := interfaces.NewInitialSetupInterface(repo, repo)
	rateMovieInterface := interfaces.NewRateMovieInterface(repo, repo)

	router.Get("/genres", initialSetupInterface.GetAllGenres)
	router.Route("/users", func(r chi.Router) {
		// r.Get("/{userId}", initialSetupInterface.GetUserById)
		r.Post("/genres", initialSetupInterface.SaveGenrePreferences)
		r.Post("/ratings", rateMovieInterface.RateMovie)
	})
}
