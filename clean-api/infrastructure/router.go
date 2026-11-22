package infrastructure

import (
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOriginsFromEnv(),
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	return r
}

const (
	defaultAllowedOrigin = "http://localhost:3000"
	allowedOriginsEnv    = "API_ALLOWED_ORIGINS"
)

// allowedOriginsFromEnv reads the comma-separated cross-origin allowlist.
// An unset or blank value falls back to the local client origin.
func allowedOriginsFromEnv() []string {
	raw := os.Getenv(allowedOriginsEnv)
	if strings.TrimSpace(raw) == "" {
		return []string{defaultAllowedOrigin}
	}

	origins := []string{}
	for _, origin := range strings.Split(raw, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return []string{defaultAllowedOrigin}
	}

	return origins
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
