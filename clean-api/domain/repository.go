package domain

type Repository interface {
	// movies
	GetAllGenres() ([]Genre, error)
	GetMovieById(userId string) (Movie, error)
	GetMoviesDetails(movieIds []string) (MoviesDetails, error)

	// user
	GetUserById(userId string) (User, error)
	RegisterNewUser(user User) error
	SaveGenrePreferences(userId string, genres []Genre) error
	RateMovie(userId, movieId string, rating int) error
	GetGenrePreferencesCount(userId string) (uint, error)
	GetUserRatingsCount(userId string) (uint, error)

	// recommendation
	GetAllRatingsForMoviesInGenre(userId string) (PopularMovies, error)
	GetSimilairUsersAndTheirAvgRating(userId string) (UsersToCompare, error)
	GetRatedMoviesForUsersYetToBeConsidered(userId string, userIds []string) (ScoringMovies, error)
	GetAllLikedMovies(userId string) (UsersLikedMovies, error)
	GetMoviesCasts(movieIds []string, movies MoviesWithoutCastDetails) error
	GetSimilarMoviesToAlreadyLikedOnes(userId string, movieIds []string) (SimilarMoviesToLikedOnes, error)

	// watchlist
	AddToWatchlist(userId, movieId string, timeOfAdding int64) error
	RemoveFromWatchlist(userId, movieId string) error
	GetWatchlist(userId string) (UserWatchlist, error)
	GetWatchlistHistory(userId string) (UserWatchlist, error)
}
