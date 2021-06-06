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

func (cbrs ContentBasedRecommendationService) GetContentBasedRecommendation(userId string, limit uint) (domain.ContentBasedRecommendations, error) {
	emptyRec := domain.ContentBasedRecommendations{}

	if _, err := cbrs.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return emptyRec, errStack
	}

	likedMovies, err := cbrs.repo.GetAllLikedMovies(userId)
	if err != nil {
		errStack := errors.Wrap(err, noLikedMovies)
		return emptyRec, errStack
	}

	if len(likedMovies) < 1 {
		cause := errors.New("not_found")
		return emptyRec, errors.Wrap(cause, noLikedMovies)
	}

	likedMovieIds := getIdsFromLikedMovies(&likedMovies)
	if err := cbrs.repo.GetMoviesCasts(likedMovieIds, likedMovies); err != nil {
		errStack := errors.Wrap(err, "could not get movie casts for likedMovies")
		return emptyRec, errStack
	}

	similarMovies, err := cbrs.repo.GetSimilarMoviesToAlreadyLikedOnes(userId, likedMovieIds)
	if err != nil {
		errStack := errors.Wrap(err, noSimilarMovies)
		return emptyRec, errStack
	}

	if len(similarMovies) < 1 {
		cause := errors.New("not_found")
		return emptyRec, errors.Wrap(cause, noSimilarMovies)
	}

	similarMoviesIds := getIdsFromSimilarMovies(&similarMovies)
	if err := cbrs.repo.GetMoviesCasts(similarMoviesIds, similarMovies); err != nil {
		errStack := errors.Wrap(err, "could not get movie casts for similarMovies")
		return emptyRec, errStack
	}

	recs, err := domain.CalculateJaccard(likedMovies, similarMovies)
	if err != nil {
		return emptyRec, err
	}

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

const (
	noLikedMovies   = "you have not rated any movies yet, please return and complete the setup"
	noSimilarMovies = "could not find similar movies to recommend, please rate some more and try again"
)
