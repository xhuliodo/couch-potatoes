package interfaces

import (
	"fmt"
	"net/http"

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
	Poster       string `json:"posterUrl"`
	MoreInfoLink string `json:"moreInfoLink"`
}

type popularMoviesView struct {
	Movie        movieView `json:"movie"`
	AvgRating    float32   `json:"avgRating"`
	CountRatings uint      `json:"ratingsCount"`
}

func (pmr popularMoviesResource) GetPopularMoviesBasedOnGenre(w http.ResponseWriter, r *http.Request) {
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	limitInterface := r.Context().Value("limit")
	limit := limitInterface.(uint)

	skipInterface := r.Context().Value("skip")
	skip := skipInterface.(uint)

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
				Poster:       movie.Poster,
				MoreInfoLink: movie.MoreInfoLink,
			},
			AvgRating:    movie.AvgRating,
			CountRatings: movie.CountRatings,
		})
	}

	render.Respond(w, r, view)
}
