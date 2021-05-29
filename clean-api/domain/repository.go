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
	GetGenrePreferences(userId string) ([]Genre, error)
	GetUserRatingsCount(userId string) (uint, error)

	// recommendation
	GetAllRatingsForMoviesInGenre(userId string) (PopularMovies, error)
	GetSimilairUsersAndTheirAvgRating(userId string) (UsersToCompare, error)
	GetRatedMoviesForUsersYetToBeConsidered(userId string, userIds []string) (ScoringMovies, error)
}
