package infrastructure

import (
	"flag"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/xhuliodo/couch-potatoes/clean-api/docs"

	"github.com/xhuliodo/couch-potatoes/clean-api/infrastructure/auth"
	"github.com/xhuliodo/couch-potatoes/clean-api/interfaces"
)

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return r
}

var routes = flag.Bool("routes", false, "Generate router documentation")

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

	router.Mount("/swagger", httpSwagger.WrapHandler)

	// opts := docsMiddleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// sh := docsMiddleware.Redoc(opts, nil)

	// router.Handle("/docs", sh)
	// router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// if *routes {
	// 	// fmt.Println(docgen.JSONRoutesDoc(r))
	// 	fmt.Println(docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{
	// 		ProjectPath: "github.com/go-chi/chi/v5",
	// 		Intro:       "Welcome to the chi/_examples/rest generated docs.",
	// 	}))
	// 	return
	// }

}
