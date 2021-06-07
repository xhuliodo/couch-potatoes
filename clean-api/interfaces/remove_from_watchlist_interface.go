package interfaces

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type removeFromWatchlistResource struct {
	removeFromWatchlistService application.RemoveFromWatchlistService
}

func NewRemoveFromWatchlistResource(removeFromWatchlistService application.RemoveFromWatchlistService) removeFromWatchlistResource {
	return removeFromWatchlistResource{removeFromWatchlistService}
}

// @router /watchlist/{movieId} [delete]
// @param movieId path int true "the id of the movie the user is removing to their watchlist"
// @param authorization header string true "Bearer token"
// @summary remove a movie from a user's watchlist or watchlist history
// @tags watchlists
// @produce  json
// @success 204 {object} EmptyView "no response it it's successful"
// @failure 404 {object} common_http.ErrorResponse "when either a movie or a user does not exist"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (rfwr removeFromWatchlistResource) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	movieId := chi.URLParam(r, "movieId")

	if errStack := rfwr.removeFromWatchlistService.RemoveFromWatchlist(userId, movieId); errStack != nil {
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	render.Render(w, r, common_http.NoContent())
}
