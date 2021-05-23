package domain

type Repository interface {
	// movies
	GetAllGenres() ([]Genre, error)
	GetMovieById(userId string) (Movie, error)
	GetAllRatingsForMoviesInGenre(userId string, genres []Genre) ([]AggregateMovieRatings, error)

	// user
	GetUserById(userId string) (User, error)
	RegisterNewUser(user User) error
	SaveGenrePreferences(userId string, genres []Genre) error
	RateMovie(userId, movieId string, rating int) error
	GetGenrePreferences(userId string) ([]Genre, error)
	// TODO: implementing
	GetUserRatingsCount(userId string) (uint, error)
}
