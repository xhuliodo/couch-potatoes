package application

import (
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type PopularMovieService struct {
	repo domain.Repository
}

func NewPopularMovieService(repo domain.Repository) PopularMovieService {
	return PopularMovieService{repo}
}

func (pms PopularMovieService) GetPopularMoviesBasedOnGenre(userId string, limit uint, skip uint) (domain.PopularMovies, error) {
	emptyRec := []domain.PopularMovie{}
	if _, err := pms.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return emptyRec, errStack
	}

	movies, err := pms.repo.GetAllRatingsForMoviesInGenre(userId)
	if err != nil {
		return emptyRec, errors.Wrap(err, noRatingsForMoviesInGenreInDb)
	}

	if len(movies) < 1 {
		cause := errors.New("not_found")
		return emptyRec, errors.Wrap(cause, noMoviesInPreferedGenres)
	}

	movies.CalculateBoostedScore()

	movies.SortMoviesBasedOnRatings()

	// handle pagination
	length := len(movies)
	begin, end, err := handlePagination(uint(length), skip, limit)
	if err != nil {
		return emptyRec, err
	}
	recs := movies[begin:end]

	moviesIds := getMovieIds(recs)
	moviesDetails, _ := pms.repo.GetMoviesDetails(moviesIds)

	recs.PopulateMoviesWithDetails(moviesDetails)

	return recs, nil
}

func getMovieIds(recs domain.PopularMovies) []string {
	moviesIds := []string{}
	for _, rec := range recs {
		movieId := rec.Movie.Id
		moviesIds = append(moviesIds, movieId)
	}
	return moviesIds
}

const (
	noRatingsForMoviesInGenreInDb = "could not get ratings for movies in genre"
	noMoviesInPreferedGenres      = "there are no movies with ratings in the prefered genres"
	tooMuchSkipping               = "you're all caught up"
)
