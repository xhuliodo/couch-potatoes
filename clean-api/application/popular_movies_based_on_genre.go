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
	sortedList, _ := SortMoviesBasedOnRatings(moviesWithRating)

	return sortedList[skip : skip+limit], nil
}

func SortMoviesBasedOnRatings(aggregateRatings []domain.AggregateMovieRatings) ([]domain.PopulatiryScoredMovie, int) {
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
		}
		unsortedList = append(unsortedList, entry)
	}

	sort.SliceStable(unsortedList, func(i, j int) bool {
		if unsortedList[i].CountRatings != unsortedList[j].CountRatings {
			return unsortedList[i].CountRatings > unsortedList[j].CountRatings
		}
		return unsortedList[i].AvgRating > unsortedList[j].AvgRating
	})

	length := len(unsortedList)

	return unsortedList, length
}

// func End(begin, end uint, len int) (begin, end uint) {

// }
