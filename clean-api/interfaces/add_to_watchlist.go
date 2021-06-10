package interfaces

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type addToWatchlistResource struct {
	addToWatchlistService application.AddToWatchlistService
}

func NewAddToWatchlistResource(addToWatchlistService application.AddToWatchlistService) addToWatchlistResource {
	return addToWatchlistResource{addToWatchlistService}
}

// @router /watchlist/{movieId} [post]
// @param movieId path int true "the id of the movie the user is adding to their watchlist"
// @param authorization header string true "Bearer token"
// @summary add a movie to a user's watchlist
// @tags watchlists
// @produce  json
// @success 201 {object} common_http.InfoResponse "api response"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 404 {object} common_http.ErrorResponse "when either a movie or a user does not exist"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (atwr addToWatchlistResource) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	movieId := chi.URLParam(r, "movieId")

	if errStack := atwr.addToWatchlistService.AddToWatchlist(userId, movieId); errStack != nil {
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	successMsg := fmt.Sprintf("you just added to watchlist the movie with id: %s", movieId)
	render.Render(w, r, common_http.ResourceCreated(successMsg))
}
