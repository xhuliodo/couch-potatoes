package infrastructure

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type InMemoryRepository struct {
	movies map[string]domain.Movie
	users  map[string]domain.User
}

func NewInMemoryRepository() *InMemoryRepository {
	ids := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25"}
	inMemRepo := InMemoryRepository{movies: map[string]domain.Movie{}, users: map[string]domain.User{}}
	// populate mock movies
	inMemRepo.movies[ids[0]] = domain.Movie{
		Title:        "Toy Story",
		ReleaseYear:  2000,
		MoreInfoLink: "nana",
		PeopleInvolved: []domain.Cast{
			{Id: ids[1]},
			{Id: ids[2]},
		},
	}
	inMemRepo.movies[ids[3]] = domain.Movie{
		Title:        "Social Network",
		ReleaseYear:  2000,
		MoreInfoLink: "nana",
		PeopleInvolved: []domain.Cast{
			{Id: ids[7]},
			{Id: ids[4]},
		},
	}
	inMemRepo.movies[ids[5]] = domain.Movie{
		Title:        "Nocturnal Animals",
		ReleaseYear:  2000,
		MoreInfoLink: "nana",
		PeopleInvolved: []domain.Cast{
			{Id: ids[1]},
			{Id: ids[2]},
		},
	}

	inMemRepo.movies[ids[6]] = domain.Movie{
		Title:        "Batman Begins",
		ReleaseYear:  2000,
		MoreInfoLink: "nana",
		PeopleInvolved: []domain.Cast{
			{Id: ids[1]},
			{Id: ids[2]},
			{Id: ids[4]},
		},
	}
	// populate users
	inMemRepo.users[ids[6]] = domain.User{
		RatedMovies: []domain.RatedMovie{
			{Movie: domain.Movie{Id: ids[0]}, Rating: 1},
			{Movie: domain.Movie{Id: ids[3]}, Rating: 1},
		},
	}

	return &inMemRepo

}

func (imr *InMemoryRepository) GetAllGenres() ([]domain.Genre, error) {
	emptyGenres := []domain.Genre{}
	return emptyGenres, nil
}

func (imr *InMemoryRepository) GetUserById(userId string) (domain.User, error) {
	emptyUser := domain.User{}
	user, found := imr.users[userId]
	if !found {
		return emptyUser, errors.New("user does not exist")
	}
	return user, nil
}

func (imr *InMemoryRepository) SaveGenrePreferences(userId string, genres []domain.Genre) error {
	for _, user := range imr.users {
		if user.Id == userId {
			user.FavoriteGenres = append(user.FavoriteGenres, genres...)
		}
	}
	return nil
}

func (imr *InMemoryRepository) GetMovieById(movieId string) (domain.Movie, error) {
	emptyMovie := domain.Movie{}
	return emptyMovie, nil
}

func (imr *InMemoryRepository) RegisterNewUser(user domain.User) error {
	return nil
}

func (imr *InMemoryRepository) GetGenrePreferences(userId string) ([]domain.Genre, error) {
	emptyGenres := []domain.Genre{}
	return emptyGenres, nil
}

func (imr *InMemoryRepository) GetAllRatingsForMoviesInGenre(userId string) (domain.PopularMovies, error) {
	emptyPopularMovies := domain.PopularMovies{}
	return emptyPopularMovies, nil
}

func (imr *InMemoryRepository) GetUserRatingsCount(userId string) (uint, error) {
	emptyCount := 0
	return uint(emptyCount), nil
}

func (imr *InMemoryRepository) GetSimilairUsersAndTheirAvgRating(userId string) (domain.UsersToCompare, error) {
	emptyUserToCompare := domain.UsersToCompare{}
	return emptyUserToCompare, nil
}

func (imr *InMemoryRepository) GetRatedMoviesForUsersYetToBeConsidered(
	userId string,
	userIds []string,
) (domain.ScoringMovies, error) {
	emptyScoringMovies := domain.ScoringMovies{}
	return emptyScoringMovies, nil
}

func (imr *InMemoryRepository) RateMovie(userId, movieId string, rating int) error {
	return nil
}

func (imr *InMemoryRepository) GetAllLikedMovies(userId string) (domain.UsersLikedMovies, error) {
	emptyLikedMovies := domain.UsersLikedMovies{}
	ratedMovies := imr.users[userId].RatedMovies
	for _, movie := range ratedMovies {
		if movie.Rating == 1 {
			emptyLikedMovies[movie.Id] = domain.UsersLikedMovie{AllCast: map[string]bool{}}
		}
	}

	return emptyLikedMovies, nil
}
func (imr *InMemoryRepository) GetMoviesCasts(movieIds []string, movies domain.MoviesWithoutCastDetails) error {
	for _, movieId := range movieIds {
		m, found := imr.movies[movieId]
		if !found {
			return errors.New("a movie with this identifier does not exist")
		}
		castSlice := []string{}
		for _, cast := range m.PeopleInvolved {
			castSlice = append(castSlice, cast.Id)
		}
		movies.PopulateWithCast(movieId, castSlice)
	}
	return nil
}
func (imr *InMemoryRepository) GetSimilarMoviesToAlreadyLikedOnes(userId string, movieIds []string) (domain.SimilarMoviesToLikedOnes, error) {
	similarMovie := domain.SimilarMoviesToLikedOnes{}
	similarMovie["6"] = domain.SimilarMovieToLikedOnes{}
	similarMovie["7"] = domain.SimilarMovieToLikedOnes{}

	return similarMovie, nil
}

func (imr *InMemoryRepository) GetMoviesDetails(movieIds []string) (domain.MoviesDetails, error) {
	emptyMovieDetails := domain.MoviesDetails{}
	for _, movieId := range movieIds {
		m, found := imr.movies[movieId]
		if !found {
			return emptyMovieDetails, errors.New("a movie with this identifier does not exist")
		}
		emptyMovieDetails[movieId] = domain.MovieDetails{Title: m.Title, ReleaseYear: m.ReleaseYear, MoreInfoLink: m.MoreInfoLink}
	}
	return emptyMovieDetails, nil
}

func (imr *InMemoryRepository) AddToWatchlist(userId, movieId string, timeOfAdding int64) error {
	return nil
}

func (imr *InMemoryRepository) RemoveFromWatchlist(userId, movieId string) error {
	return nil
}
