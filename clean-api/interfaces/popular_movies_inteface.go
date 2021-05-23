package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type popularMoviesResource struct {
	popularMoviesService application.PopularMovieService
}

type movieView struct {
	Id           string `json:"genreId"`
	Title        string `json:"title"`
	ReleaseYear  int    `json:"releaseYear"`
	MoreInfoLink string `json:"moreInfoLink"`
}

type popularMoviesView struct {
	Movie        movieView `json:"movie"`
	AvgRating    float32   `json:"avgRating"`
	CountRatings uint      `json:"ratingsCount"`
	GenreMatched uint      `json:"genreMatched"`
}

func (pmr popularMoviesResource) GetPopularMoviesBasedOnGenre(w http.ResponseWriter, r *http.Request) {
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	var limit uint
	limitUrlQueryParam := r.URL.Query().Get("limit")
	limitU64, err := strconv.ParseUint(limitUrlQueryParam, 10, 32)
	if err != nil {
		limit = 5
	}
	limit = uint(limitU64)

	var skip uint
	skipUrlQueryParam := r.URL.Query().Get("skip")
	skipU64, err := strconv.ParseUint(skipUrlQueryParam, 10, 32)
	if err != nil {
		skip = 5
	}
	skip = uint(skipU64)

	popularMovies, err := pmr.popularMoviesService.GetPopularMoviesBasedOnGenre(userId, limit, skip)
	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	view := []popularMoviesView{}
	for _, movie := range popularMovies {
		view = append(view, popularMoviesView{
			Movie: movieView{
				Id:           string(movie.Id),
				Title:        movie.Title,
				ReleaseYear:  movie.ReleaseYear,
				MoreInfoLink: movie.MoreInfoLink,
			},
			AvgRating:    movie.AvgRating,
			CountRatings: movie.CountRatings,
			GenreMatched: movie.GenreMatched,
		})
	}

	render.Respond(w, r, view)
}
