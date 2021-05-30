package application

import (
	"errors"

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
		return emptyRec, errors.New("a user with this identifier does not exist")
	}

	movies, err := pms.repo.GetAllRatingsForMoviesInGenre(userId)
	if err != nil {
		return emptyRec, errors.New("you're all caught up")
	}

	movies.CalculateBoostedScore()

	movies.SortMoviesBasedOnRatings()

	// handle pagination
	length := len(movies)
	begin, end, err := handlePagination(uint(length), skip, limit)
	if err != nil {
		return emptyRec, errors.New("you're all caught up")
	}
	recs := movies[begin:end]

	moviesIds := getMovieIds(recs)
	moviesDetails, _ := pms.repo.GetMoviesDetails(moviesIds)

	recs.PopulateMoviesWithDetails(moviesDetails)

	return recs, nil
}

const (
	defaultLimit uint = 5
	defaultSkip  uint = 0
)

func handlePagination(len, skip, limit uint) (begin, end uint, err error) {
	if skip != defaultSkip {
		maxSkip := maxSkip(uint(len), limit)
		if skip > maxSkip {
			return 0, 0, errors.New("you've reached the limit")
		}
	}

	begin = skip

	remaining := len - skip
	if remaining < limit {
		end = remaining + skip
		return begin, end, nil
	}

	end = limit + skip
	return begin, end, nil

}

func maxSkip(total uint, limit uint) uint {
	if limit == 0 {
		limit = defaultLimit
	}
	return ((total - 1) / limit) * limit
}

func getMovieIds(recs domain.PopularMovies) []string {
	moviesIds := []string{}
	for _, rec := range recs {
		movieId := rec.Movie.Id
		moviesIds = append(moviesIds, movieId)
	}
	return moviesIds
}
