package domain

import (
	"sort"
)

type PopularMovies []PopularMovie

type PopularMovie struct {
	Movie        Movie
	AvgRating    float64
	RatingsCount uint
	GenreMatched uint
	BoostedScore uint
}

type MoviesDetails map[string]MovieDetails

type MovieDetails struct {
	Title        string
	ReleaseYear  int
	MoreInfoLink string
}

func (pm *PopularMovies) CalculateBoostedScore() {
	for i, movie := range *pm {
		countRatings := movie.RatingsCount
		genreMultiplier := movie.GenreMatched
		scoreBoosted := countRatings * genreMultiplier

		(*pm)[i].BoostedScore = scoreBoosted
	}
}

func (pm *PopularMovies) SortMoviesBasedOnRatings() {
	sort.SliceStable(*pm, func(i, j int) bool {
		if (*pm)[i].BoostedScore != (*pm)[j].BoostedScore {
			return (*pm)[i].BoostedScore > (*pm)[j].BoostedScore
		}
		return (*pm)[i].AvgRating > (*pm)[j].AvgRating
	})
}

func (pm *PopularMovies) PopulateMoviesWithDetails(moviesDetails MoviesDetails) {
	for i, movie := range *pm {
		(*pm)[i].Movie.Title = moviesDetails[movie.Movie.Id].Title
		(*pm)[i].Movie.ReleaseYear = moviesDetails[movie.Movie.Id].ReleaseYear
		(*pm)[i].Movie.MoreInfoLink = moviesDetails[movie.Movie.Id].MoreInfoLink
	}
}
