package domain

import (
	"sort"
)

type MoviesWithoutCastDetails interface {
	PopulateWithCast(movieId string, castDetails []string)
}

type UsersLikedMovies map[string]UsersLikedMovie

type UsersLikedMovie struct {
	AllCast map[string]bool
}

func (ulm UsersLikedMovies) PopulateWithCast(movieId string, castDetails []string) {
	for _, castId := range castDetails {
		ulm[movieId].AllCast[castId] = true

	}
}

type SimilarMoviesToLikedOnes map[string]SimilarMovieToLikedOnes

type SimilarMovieToLikedOnes struct {
	AllCast []string
}

func (smtlo SimilarMoviesToLikedOnes) PopulateWithCast(movieId string, castDetails []string) {
	smtlo[movieId] = SimilarMovieToLikedOnes{
		AllCast: castDetails,
	}

}

type ContentBasedRecommendations []ContentBasedRecommendation

type ContentBasedRecommendation struct {
	Movie
	Score float64
}

func (cbr *ContentBasedRecommendations) SortByScoreDesc() {
	sort.SliceStable(*cbr, func(i, j int) bool {
		return (*cbr)[i].Score > (*cbr)[j].Score
	})
}

func (cbr ContentBasedRecommendations) PopulateWithMovieDetails(moviesDetails MoviesDetails) {
	for i, movie := range cbr {
		cbr[i].Movie.Title = moviesDetails[movie.Movie.Id].Title
		cbr[i].Movie.ReleaseYear = moviesDetails[movie.Movie.Id].ReleaseYear
		cbr[i].Movie.MoreInfoLink = moviesDetails[movie.Movie.Id].MoreInfoLink
	}
}

const leastThingsInCommon float64 = 2

func CalculateJaccard(ulm UsersLikedMovies, similarMovies SimilarMoviesToLikedOnes) ContentBasedRecommendations {
	recs := ContentBasedRecommendations{}

	for _, likedMovie := range ulm {
		for similarMovieId, similarMovie := range similarMovies {
			intersection := []string{}
			union := similarMovie.AllCast
			for _, castId := range similarMovie.AllCast {
				if ok := likedMovie.AllCast[castId]; !ok {
					union = append(union, castId)
				} else {
					intersection = append(intersection, castId)
				}
			}
			unionCount := float64(len(union))
			intersectionCount := float64(len(intersection))
			var jaccard float64
			if intersectionCount >= leastThingsInCommon {
				jaccard = intersectionCount / unionCount
				rec := ContentBasedRecommendation{
					Movie: Movie{Id: similarMovieId},
					Score: jaccard,
				}

				recs = append(recs, rec)
			}
		}
	}
	return recs
}

func (cbr ContentBasedRecommendations) RemoveDuplicates() ContentBasedRecommendations {
	indexList := make(map[string]float64)

	for _, movie := range cbr {
		if indexList[movie.Id] < movie.Score {
			indexList[movie.Id] = movie.Score
		}
	}

	recsWithNoDups := ContentBasedRecommendations{}
	for movieId, jaccard := range indexList {
		rec := ContentBasedRecommendation{
			Movie: Movie{Id: movieId},
			Score: jaccard,
		}
		recsWithNoDups = append(recsWithNoDups, rec)
	}
	return recsWithNoDups
}
