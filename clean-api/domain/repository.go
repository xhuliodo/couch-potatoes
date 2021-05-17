package domain

type Repository interface {
	// movies
	GetAllGenres() ([]Genre, error)
	GetMovieById(userId string) (Movie, error)
	// TODO: implementing
	GetMoviesInGenre(genres []Genre) ([]Movie, error)
	GetAllRatingsForMovies(movie []Movie) ([]AggregateMovieRatings, error)

	// user
	GetUserById(userId string) (User, error)
	RegisterNewUser(user User) error
	SaveGenrePreferences(userId string, genres []Genre) error
	RateMovie(userId, movieId string, rating int) error
	// TODO: implementing
	GetGenrePreferences(userId string) ([]Genre, error)
}
