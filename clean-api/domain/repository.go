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
	GetUserRatingsCount(userId string) (uint, error)
	GetAvgRatingAndCollectSimilairUsers(userId string) (UserComparison, error)
	GetUsersAvgRating(userIds []string, userComparison *UserComparison) error
	GetRatedMoviesForUsers(userIds []string) ([]User, error)

	// recommendation
}
