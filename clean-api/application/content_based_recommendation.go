package application

import (
	"github.com/pkg/errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type ContentBasedRecommendationService struct {
	repo domain.Repository
}

func NewContentBasedRecommendationService(repo domain.Repository) ContentBasedRecommendationService {
	return ContentBasedRecommendationService{repo}
}

const (
	noLikedMovies   = "you have not rated any movies yet, please return and complete the setup"
	noSimilarMovies = "could not find similar movies to recommend, please rate some more and try again"
)

func (cbrs ContentBasedRecommendationService) GetContentBasedRecommendation(userId string, limit uint) (
	recs domain.ContentBasedRecommendations, err error,
) {
	if _, err := cbrs.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return recs, errStack
	}

	likedMovies, err := cbrs.repo.GetAllLikedMovies(userId)
	if err != nil {
		errStack := errors.Wrap(err, noLikedMovies)
		return recs, errStack
	}

	if len(likedMovies) < 1 {
		cause := errors.New("not_found")
		return recs, errors.Wrap(cause, noLikedMovies)
	}

	likedMovieIds := getIdsFromLikedMovies(&likedMovies)
	if err := cbrs.repo.GetMoviesCasts(likedMovieIds, likedMovies); err != nil {
		errStack := errors.Wrap(err, "could not get movie casts for likedMovies")
		return recs, errStack
	}

	similarMovies, err := cbrs.repo.GetSimilarMoviesToAlreadyLikedOnes(userId, likedMovieIds)
	if err != nil {
		errStack := errors.Wrap(err, noSimilarMovies)
		return recs, errStack
	}

	if len(similarMovies) < 1 {
		cause := errors.New("not_found")
		return recs, errors.Wrap(cause, noSimilarMovies)
	}

	similarMoviesIds := getIdsFromSimilarMovies(&similarMovies)
	if err := cbrs.repo.GetMoviesCasts(similarMoviesIds, similarMovies); err != nil {
		errStack := errors.Wrap(err, "could not get movie casts for similarMovies")
		return recs, errStack
	}

	recs, err = domain.CalculateJaccard(likedMovies, similarMovies)
	if err != nil {
		return recs, err
	}

	recsWithNoDups := recs.RemoveDuplicates()

	recsWithNoDups.SortByScoreDesc()

	// handle pagination
	length := len(recsWithNoDups)
	begin, end, err := handlePagination(uint(length), defaultSkip, limit)
	if err != nil {
		return recs, err
	}
	recs = recsWithNoDups[begin:end]

	remainingRecsIds := getIdFromRemainingRecs(recs)
	moviesDetails, _ := cbrs.repo.GetMoviesDetails(remainingRecsIds)
	recs.PopulateWithMovieDetails(moviesDetails)

	return recs, nil
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
