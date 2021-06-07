package interfaces

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type popularMoviesResource struct {
	popularMoviesService application.PopularMovieService
}

type movieView struct {
	Id           string `json:"movieId"`
	Title        string `json:"title"`
	ReleaseYear  int    `json:"releaseYear"`
	MoreInfoLink string `json:"moreInfoLink"`
}//@name MovieResponse

type popularMoviesView struct {
	Movie        movieView `json:"movie"`
	BoostedScore uint      `json:"score"`
}//@name PopularMoviesRecommendationResponse

// @router /recommendations/popular [get]
// @param skip query int true "skip" minimum(0) default(0)
// @param limit query int true "limit" default(5)
// @param authorization header string true "Bearer token"
// @summary get most popular movies based on provided genre preferences
// @tags recommendations
// @produce json
// @success 200 {object} common_http.SuccessResponse{data=popularMoviesView} "api response"
// @failure 400 {object} common_http.ErrorResponse "when the skip query param gets too big"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (pmr popularMoviesResource) GetPopularMoviesBasedOnGenre(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	limit := getLimit(r)
	skip := getSkip(r)

	popularMovies, err := pmr.popularMoviesService.GetPopularMoviesBasedOnGenre(userId, limit, skip)
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	view := []popularMoviesView{}
	for _, movie := range popularMovies {
		view = append(view, popularMoviesView{
			Movie: movieView{
				Id:           string(movie.Movie.Id),
				Title:        movie.Movie.Title,
				ReleaseYear:  movie.Movie.ReleaseYear,
				MoreInfoLink: movie.Movie.MoreInfoLink,
			},
			BoostedScore: movie.BoostedScore,
		})
	}

	render.Render(w, r, common_http.SendPayload(view))
}
