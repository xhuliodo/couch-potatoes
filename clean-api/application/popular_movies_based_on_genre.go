package application

import (
	"errors"
	"sort"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type PopularMovieService struct {
	repo domain.Repository
}

func NewPopularMovieService(repo domain.Repository) PopularMovieService {
	return PopularMovieService{repo}
}

func (pms PopularMovieService) GetPopularMoviesBasedOnGenre(userId string, limit uint, skip uint) ([]domain.PopulatiryScoredMovie, error) {
	emptyRec := []domain.PopulatiryScoredMovie{}
	if _, err := pms.repo.GetUserById(userId); err != nil {
		return emptyRec, errors.New("a user with this identifier does not exist")
	}

	genres, err := pms.repo.GetGenrePreferences(userId)
	if err != nil {
		return emptyRec, err
	}

	moviesWithRating, err := pms.repo.GetAllRatingsForMoviesInGenre(userId, genres)
	if err != nil {
		return emptyRec, errors.New("you're all caught up")
	}

	// TODO: handle weirdly small recommendations
	sortedList, length := SortMoviesBasedOnRatings(moviesWithRating)

	begin, end, err := handlePagination(length, skip, limit)

	if err != nil {
		return emptyRec, errors.New("you're all caught up")
	}

	return sortedList[begin:end], nil
}

func SortMoviesBasedOnRatings(aggregateRatings []domain.AggregateMovieRatings) (sortedList []domain.PopulatiryScoredMovie, length uint) {
	unsortedList := []domain.PopulatiryScoredMovie{}

	for _, movieWithAllRatings := range aggregateRatings {
		ratings := movieWithAllRatings.AllRatings
		var ratingSum float32 = 0.0
		for _, rating := range ratings {
			ratingSum += rating
		}

		countRatings := len(ratings)
		avgRating := ratingSum / float32(countRatings)

		genreMultiplier := movieWithAllRatings.GenreMatched
		countBoosted := countRatings * int(genreMultiplier)
		entry := domain.PopulatiryScoredMovie{
			Movie:        movieWithAllRatings.Movie,
			CountRatings: uint(countBoosted),
			AvgRating:    avgRating,
			GenreMatched: genreMultiplier,
		}
		unsortedList = append(unsortedList, entry)
	}

	sort.SliceStable(unsortedList, func(i, j int) bool {
		if unsortedList[i].CountRatings != unsortedList[j].CountRatings {
			return unsortedList[i].CountRatings > unsortedList[j].CountRatings
		}
		return unsortedList[i].AvgRating > unsortedList[j].AvgRating
	})

	len := len(unsortedList)

	return unsortedList, uint(len)
}

const (
	defaultLimit uint = 5
	defaultSkip  uint = 0
)

func handlePagination(len, skip, limit uint) (begin, end uint, err error) {
	if skip != defaultSkip {
		maxSkip := maxSkip(len, limit)
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
