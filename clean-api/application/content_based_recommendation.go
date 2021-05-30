package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type ContentBasedRecommendationService struct {
	repo domain.Repository
}

func NewContentBasedRecommendationService(repo domain.Repository) ContentBasedRecommendationService {
	return ContentBasedRecommendationService{repo}
}

func (cbrs ContentBasedRecommendationService) GetContentBasedRecommendation(userId string, limit uint) (domain.ContentBasedRecommendations, error) {
	emptyRec := domain.ContentBasedRecommendations{}

	if _, err := cbrs.repo.GetUserById(userId); err != nil {
		return emptyRec, errors.New("a user with this identifier does not exist")
	}

	likedMovies, err := cbrs.repo.GetAllLikedMovies(userId)
	if err != nil {
		return emptyRec, errors.New("you have not rated any movies yet, please return and complete the setup")
	}

	likedMovieIds := getIdsFromLikedMovies(&likedMovies)
	if err := cbrs.repo.GetMoviesCasts(likedMovieIds, likedMovies); err != nil {
		return emptyRec, err
	}

	similarMovies, err := cbrs.repo.GetSimilarMoviesToAlreadyLikedOnes(userId, likedMovieIds)
	if err != nil {
		return emptyRec, err
	}

	similarMoviesIds := getIdsFromSimilarMovies(&similarMovies)
	if err := cbrs.repo.GetMoviesCasts(similarMoviesIds, similarMovies); err != nil {
		return emptyRec, err
	}

	recs := domain.CalculateJaccard(likedMovies, similarMovies)

	recsWithNoDups := recs.RemoveDuplicates()

	recsWithNoDups.SortByScoreDesc()

	// handle pagination
	length := len(recsWithNoDups)
	begin, end, err := handlePagination(uint(length), defaultSkip, limit)
	if err != nil {
		return emptyRec, errors.New("you're all caught up")
	}
	remainingRecs := recsWithNoDups[begin:end]

	remainingRecsIds := getIdFromRemainingRecs(remainingRecs)

	moviesDetails, _ := cbrs.repo.GetMoviesDetails(remainingRecsIds)

	remainingRecs.PopulateWithMovieDetails(moviesDetails)

	return remainingRecs, nil
}

func getIdsFromLikedMovies(movies *domain.UsersLikedMovies) []string {
	movieIds := []string{}
	for key := range *movies {
		movieIds = append(movieIds, key)
	}
	return movieIds
}

func getIdsFromSimilarMovies(movies *domain.SimilarMoviesToLikedOnes) []string {
	movieIds := []string{}
	for key := range *movies {
		movieIds = append(movieIds, key)
	}
	return movieIds
}

func getIdFromRemainingRecs(recs domain.ContentBasedRecommendations) []string {
	movieIds := []string{}
	for _, movie := range recs {
		movieIds = append(movieIds, movie.Id)
	}
	return movieIds
}
